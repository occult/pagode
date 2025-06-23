package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/user"
	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/log"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/redirect"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	"github.com/occult/pagode/pkg/ui"
	"github.com/resend/resend-go/v2"

	inertia "github.com/romsar/gonertia/v2"
)

type Auth struct {
	config  *config.Config
	auth    *services.AuthClient
	mail    *services.MailClient
	orm     *ent.Client
	Inertia *inertia.Inertia
}

type RegisterForm struct {
	Name            string `form:"name" validate:"required"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required"`
	ConfirmPassword string `form:"password_confirmation" validate:"required,eqfield=Password"`
	form.Submission
}

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
	form.Submission
}

type ForgotPassword struct {
	Email string `form:"email" validate:"required,email"`
	form.Submission
}

type ResetPassword struct {
	Password        string `form:"password" validate:"required"`
	ConfirmPassword string `form:"password_confirmation" validate:"required,eqfield=Password"`

	form.Submission
}

func init() {
	Register(new(Auth))
}

func (h *Auth) Init(c *services.Container) error {
	h.config = c.Config
	h.orm = c.ORM
	h.auth = c.Auth
	h.mail = c.Mail
	h.Inertia = c.Inertia
	return nil
}

func (h *Auth) Routes(g *echo.Group) {
	g.GET("/logout", h.Logout, middleware.RequireAuthentication).Name = routenames.Logout
	g.GET("/email/verify/:token", h.VerifyEmail).Name = routenames.VerifyEmail

	noAuth := g.Group("/user", middleware.RequireNoAuthentication)
	noAuth.GET("/login", h.LoginPage).Name = routenames.Login
	noAuth.POST("/login", h.LoginSubmit).Name = routenames.LoginSubmit
	noAuth.GET("/register", h.RegisterPage).Name = routenames.Register
	noAuth.POST("/register", h.RegisterSubmit).Name = routenames.RegisterSubmit
	noAuth.GET("/password", h.ForgotPasswordPage).Name = routenames.ForgotPassword
	noAuth.POST("/password", h.ForgotPasswordSubmit).Name = routenames.ForgotPasswordSubmit

	resetGroup := noAuth.Group("/password/reset",
		middleware.LoadUser(h.orm),
		middleware.LoadValidPasswordToken(h.auth),
	)
	resetGroup.GET("/token/:user/:password_token/:token", h.ResetPasswordPage).Name = routenames.ResetPassword
	resetGroup.POST("/token/:user/:password_token/:token", h.ResetPasswordSubmit).Name = routenames.ResetPasswordSubmit
}

func (h *Auth) LoginPage(ctx echo.Context) error {
	canResetPassword := true

	err := h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Auth/Login",
		inertia.Props{
			"canResetPassword": canResetPassword,
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

func (h *Auth) LoginSubmit(ctx echo.Context) error {
	w := ctx.Response().Writer
	r := ctx.Request()

	var input LoginForm

	uriLogin := ctx.Echo().Reverse(routenames.Login)

	authFailed := func() error {
		input.SetFieldError("Email", "")
		input.SetFieldError("Password", "")
		msg.Danger(ctx, "Invalid credentials. Please try again.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil
	}

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.LoginPage(ctx)
	default:
		return err
	}

	// Attempt to load the user.
	u, err := h.orm.User.
		Query().
		Where(user.Email(strings.ToLower(input.Email))).
		Only(ctx.Request().Context())

	switch err.(type) {
	case *ent.NotFoundError:
		return authFailed()
	case nil:
	default:
		return fail(err, "error querying user during login", h.Inertia, ctx)
	}

	// Check if the password is correct.
	err = h.auth.CheckPassword(input.Password, u.Password)
	if err != nil {
		return authFailed()
	}

	// Log the user in.
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		return fail(err, "unable to log in user", h.Inertia, ctx)
	}

	uriDashboard := ctx.Echo().Reverse(routenames.Dashboard)

	msg.Success(ctx, fmt.Sprintf("Welcome back, %s. You are now logged in.", u.Name))

	h.Inertia.Redirect(w, r, uriDashboard)
	return nil
}

func (h *Auth) Logout(ctx echo.Context) error {
	if err := h.auth.Logout(ctx); err == nil {
		msg.Success(ctx, "You have been logged out successfully.")
	} else {
		msg.Danger(ctx, "An error occurred. Please try again.")
	}
	return redirect.New(ctx).
		Route(routenames.Welcome).
		Go()
}

func (h *Auth) RegisterPage(ctx echo.Context) error {
	err := h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Auth/Register",
		inertia.Props{
			"text": "Teste",
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

func (h *Auth) RegisterSubmit(ctx echo.Context) error {
	w := ctx.Response().Writer
	r := ctx.Request()

	var input RegisterForm
	err := form.Submit(ctx, &input)

	uriLogin := ctx.Echo().Reverse(routenames.Login)

	log.Ctx(ctx).Info("🔍 Register form submitted", "input", input)

	// Validate submitted form data
	switch err.(type) {
	case nil:
		log.Ctx(ctx).Info("Form submitted successfully", "data", input)

	case validator.ValidationErrors:
		msg.Danger(ctx, "Please fill in all required fields correctly.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil

	default:
		msg.Danger(ctx, "Something went wrong. Please try again.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil
	}

	// Attempt to create the user
	u, err := h.orm.User.
		Create().
		SetName(input.Name).
		SetEmail(input.Email).
		SetPassword(input.Password).
		Save(r.Context())

	switch err.(type) {
	case nil:
		log.Ctx(ctx).Info("✅ User created",
			"user_name", input.Name,
			"user_email", input.Email,
		)
	case *ent.ConstraintError:
		msg.Warning(ctx, "A user with this email address already exists. Please log in.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil
	default:
		return fail(err, "unable to create user", h.Inertia, ctx)
	}

	// Try to log the user in
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		msg.Info(ctx, "Your account has been created.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil
	}

	msg.Success(ctx, "Your account has been created. You are now logged in.")

	// Send verification email
	err = h.sendVerificationEmail(ctx, u)
	if err != nil {
		log.Ctx(ctx).Error("unable to send email verification",
			"user_id", u.ID,
			"error", err,
		)
	}

	uriDashboard := ctx.Echo().Reverse(routenames.Dashboard)

	h.Inertia.Redirect(w, r, uriDashboard)
	return nil
}

func (h *Auth) sendVerificationEmail(ctx echo.Context, usr *ent.User) error {
	existingUser, err := h.orm.User.
		Query().
		Where(user.EmailEQ(usr.Email), user.IDNEQ(usr.ID)).
		Only(ctx.Request().Context())

	if err == nil && existingUser != nil {
		return fail(fmt.Errorf("email already in use"), "this email is already associated with another account", h.Inertia, ctx)
	}

	if ent.IsNotFound(err) {
		err = nil
	} else if err != nil {
		return fail(err, "failed to check for existing user by email", h.Inertia, ctx)
	}

	client := resend.NewClient(h.config.Mail.ResendApiKey)

	token, err := h.auth.GenerateEmailVerificationToken(usr.Email)
	if err != nil {
		log.Ctx(ctx).Error("unable to generate email verification token",
			"user_id", usr.ID,
			"error", err,
		)
		return fail(err, "failed to generate verification token", h.Inertia, ctx)
	}

	url := ui.NewRequest(ctx).
		Url(routenames.VerifyEmail, token)

	subject := "Confirm your email address"
	html := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thank you for signing up. Please verify your email address by clicking the link below:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you didn’t create an account, you can ignore this email.</p>
	`, usr.Name, url)

	params := &resend.SendEmailRequest{
		From:    h.config.Mail.FromAddress,
		To:      []string{usr.Email},
		Subject: subject,
		Html:    html,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Ctx(ctx).Error("unable to send email verification token",
			"user_id", usr.ID,
			"error", err,
		)
		return fail(err, "failed to send verification email", h.Inertia, ctx)
	}

	msg.Info(ctx, "An email was sent to you to verify your email address.")
	return nil
}

func (h *Auth) VerifyEmail(ctx echo.Context) error {
	var usr *ent.User

	w := ctx.Response().Writer
	r := ctx.Request()

	uriWelcome := ctx.Echo().Reverse(routenames.Welcome)

	uriForgotPassword := ctx.Echo().Reverse(routenames.ForgotPassword)

	// Validate the token.
	token := ctx.Param("token")
	email, err := h.auth.ValidateEmailVerificationToken(token)
	if err != nil {
		msg.Warning(ctx, "The link is either invalid or has expired.")
		h.Inertia.Redirect(w, r, uriForgotPassword)
	}

	// Check if it matches the authenticated user.
	if u := ctx.Get(context.AuthenticatedUserKey); u != nil {
		authUser := u.(*ent.User)

		if authUser.Email == email {
			usr = authUser
		}
	}

	// Query to find a matching user, if needed.
	if usr == nil {
		usr, err = h.orm.User.
			Query().
			Where(user.Email(email)).
			Only(ctx.Request().Context())
		if err != nil {
			return fail(err, "query failed loading email verification token user", h.Inertia, ctx)
		}
	}

	// Verify the user, if needed.
	if !usr.Verified {
		usr, err = usr.
			Update().
			SetVerified(true).
			Save(ctx.Request().Context())
		if err != nil {
			return fail(err, "failed to set user as verified", h.Inertia, ctx)
		}
	}

	msg.Success(ctx, "Your email has been successfully verified.")

	h.Inertia.Redirect(w, r, uriWelcome)
	return nil
}

func (h *Auth) ForgotPasswordPage(ctx echo.Context) error {
	err := h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Auth/ForgotPassword",
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

func (h *Auth) ForgotPasswordSubmit(ctx echo.Context) error {
	var input ForgotPassword

	w := ctx.Response().Writer
	r := ctx.Request()

	uriForgotPassword := ctx.Echo().Reverse(routenames.ForgotPassword)

	succeed := func() error {
		form.Clear(ctx)
		msg.Success(ctx, "An email containing a link to reset your password will be sent to this address if it exists in our system.")
		h.Inertia.Redirect(w, r, uriForgotPassword)
		return nil
	}

	err := form.Submit(ctx, &input)
	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ForgotPasswordPage(ctx)
	default:
		return fail(err, "form submission error on forgot password", h.Inertia, ctx)
	}

	// Attempt to load the user.
	u, err := h.orm.User.
		Query().
		Where(user.Email(strings.ToLower(input.Email))).
		Only(ctx.Request().Context())

	switch err.(type) {
	case *ent.NotFoundError:
		// We return success without revealing the email does not exist. This prevents user enumeration.
		return succeed()
	case nil:
	default:
		return fail(err, "error querying user during forgot password", h.Inertia, ctx)
	}

	client := resend.NewClient(h.config.Mail.ResendApiKey)

	// Generate the token.
	token, pt, err := h.auth.GeneratePasswordResetToken(ctx, u.ID)
	if err != nil {
		return fail(err, "error generating password reset token", h.Inertia, ctx)
	}

	log.Ctx(ctx).Info("generated password reset token",
		"user_id", u.ID,
	)

	url := ctx.Echo().Reverse(routenames.ResetPassword, u.ID, pt.ID, token)

	subject := "Reset your password"
	html := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>To reset your password go to the link below:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>If you didn’t request a password update, you can ignore this email.</p>
	`, u.Name, h.config.App.Host+url)

	params := &resend.SendEmailRequest{
		From:    h.config.Mail.FromAddress,
		To:      []string{u.Email},
		Subject: subject,
		Html:    html,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Ctx(ctx).Error("unable to send password reset email",
			"user_id", u.ID,
			"error", err,
		)
		return fail(err, "failed to send password reset email", h.Inertia, ctx)
	}

	return succeed()
}

func (h *Auth) ResetPasswordPage(ctx echo.Context) error {
	userIDStr := ctx.Param("user")
	passwordTokenID := ctx.Param("password_token")
	token := ctx.Param("token")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	u, err := h.orm.User.
		Query().
		Where(user.IDEQ(userID)).
		Only(ctx.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return fail(err, "error loading user in ResetPasswordPage", h.Inertia, ctx)
	}

	props := map[string]any{
		"token":           token,
		"userID":          userIDStr,
		"passwordTokenID": passwordTokenID,
		"email":           u.Email,
	}

	err = h.Inertia.Render(ctx.Response().Writer, ctx.Request(), "Auth/ResetPassword", props)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

func (h *Auth) ResetPasswordSubmit(ctx echo.Context) error {
	var input ResetPassword

	err := form.Submit(ctx, &input)

	w := ctx.Response().Writer
	r := ctx.Request()

	uriLogin := ctx.Echo().Reverse(routenames.Login)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		log.Ctx(ctx).Warn("⚠️ Validation failed", "errors", err)
		msg.Danger(ctx, "Please fill in the fields correctly.")
		h.Inertia.Redirect(w, r, r.URL.Path)
		return nil
	default:
		msg.Danger(ctx, "There was a problem processing your request.")
		h.Inertia.Redirect(w, r, r.URL.Path)
		return nil
	}

	usr, ok := ctx.Get(context.UserKey).(*ent.User)
	if !ok || usr == nil {
		msg.Danger(ctx, "User not found.")
		h.Inertia.Redirect(w, r, r.URL.Path)
		return nil
	}

	_, err = usr.Update().
		SetPassword(input.Password).
		Save(ctx.Request().Context())
	if err != nil {
		msg.Danger(ctx, "Unable to update your password. Please try again.")
		h.Inertia.Redirect(w, r, r.URL.Path)
		return nil
	}

	err = h.auth.DeletePasswordTokens(ctx, usr.ID)
	if err != nil {
		msg.Danger(ctx, "Password updated, but failed to clean up tokens.")
		h.Inertia.Redirect(w, r, uriLogin)
		return nil
	}

	msg.Success(ctx, "Your password has been updated.")
	h.Inertia.Redirect(w, r, uriLogin)
	return nil
}
