package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/user"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/form"
	"github.com/mikestefanello/pagoda/pkg/log"
	"github.com/mikestefanello/pagoda/pkg/middleware"
	"github.com/mikestefanello/pagoda/pkg/msg"
	"github.com/mikestefanello/pagoda/pkg/redirect"
	"github.com/mikestefanello/pagoda/pkg/routenames"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/mikestefanello/pagoda/pkg/ui/emails"

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
		return fail(err, "error querying user during login")
	}

	// Check if the password is correct.
	err = h.auth.CheckPassword(input.Password, u.Password)
	if err != nil {
		return authFailed()
	}

	// Log the user in.
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		return fail(err, "unable to log in user")
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
		return fail(err, "unable to create user")
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
	h.sendVerificationEmail(ctx, u)

	uriDashboard := ctx.Echo().Reverse(routenames.Dashboard)

	h.Inertia.Redirect(w, r, uriDashboard)
	return nil
}

func (h *Auth) sendVerificationEmail(ctx echo.Context, usr *ent.User) {
	// Generate a token.
	token, err := h.auth.GenerateEmailVerificationToken(usr.Email)
	if err != nil {
		log.Ctx(ctx).Error("unable to generate email verification token",
			"user_id", usr.ID,
			"error", err,
		)
		return
	}

	// Send the email.
	err = h.mail.
		Compose().
		To(usr.Email).
		Subject("Confirm your email address").
		Component(emails.ConfirmEmailAddress(ctx, usr.Name, token)).
		Send(ctx)
	if err != nil {
		log.Ctx(ctx).Error("unable to send email verification link",
			"user_id", usr.ID,
			"error", err,
		)
		return
	}

	msg.Info(ctx, "An email was sent to you to verify your email address.")
}

func (h *Auth) VerifyEmail(ctx echo.Context) error {
	var usr *ent.User

	// Validate the token.
	token := ctx.Param("token")
	email, err := h.auth.ValidateEmailVerificationToken(token)
	if err != nil {
		msg.Warning(ctx, "The link is either invalid or has expired.")
		return redirect.New(ctx).
			Route(routenames.Home).
			Go()
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
			return fail(err, "query failed loading email verification token user")
		}
	}

	// Verify the user, if needed.
	if !usr.Verified {
		usr, err = usr.
			Update().
			SetVerified(true).
			Save(ctx.Request().Context())
		if err != nil {
			return fail(err, "failed to set user as verified")
		}
	}

	msg.Success(ctx, "Your email has been successfully verified.")
	return redirect.New(ctx).
		Route(routenames.Home).
		Go()
}
