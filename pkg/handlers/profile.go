package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/log"
	"github.com/mikestefanello/pagoda/pkg/msg"
	"github.com/mikestefanello/pagoda/pkg/redirect"
	"github.com/mikestefanello/pagoda/pkg/routenames"
	"github.com/mikestefanello/pagoda/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Profile struct {
	orm     *ent.Client
	Inertia *inertia.Inertia
	auth    *services.AuthClient
}

func init() {
	Register(new(Profile))
}

func (h *Profile) Init(c *services.Container) error {
	h.orm = c.ORM
	h.Inertia = c.Inertia
	h.auth = c.Auth
	return nil
}

func (h *Profile) Routes(g *echo.Group) {
	profile := g.Group("/profile")
	profile.GET("/info", h.EditPage).Name = routenames.ProfileEdit
	profile.POST("/update", h.UpdateBasicInfo).Name = routenames.ProfileUpdate
	profile.DELETE("/delete", h.Delete).Name = routenames.ProfileDestroy

	profile.GET("/appearance", h.AppearancePage).Name = routenames.ProfileAppearance
	profile.GET("/password", h.PasswordPage).Name = routenames.ProfilePassword
	profile.POST("/update-password", h.UpdatePassword).Name = routenames.ProfileUpdatePassword
}

func (h *Profile) EditPage(ctx echo.Context) error {
	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Settings/Profile",
	)
}

func (h *Profile) UpdateBasicInfo(ctx echo.Context) error {
	usr, ok := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	if !ok {
		msg.Danger(ctx, "You must be logged in.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	name := ctx.FormValue("name")
	email := ctx.FormValue("email")

	if name == "" && email == "" {
		msg.Info(ctx, "Nothing to update.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	update := h.orm.User.UpdateOne(usr)

	if name != "" {
		update = update.SetName(name)
	}
	if email != "" {
		update = update.SetEmail(email)
	}

	_, err := update.Save(ctx.Request().Context())
	if err != nil {
		msg.Danger(ctx, "Failed to update user")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	msg.Success(ctx, "Your profile has been updated.")
	h.Inertia.Back(ctx.Response().Writer, ctx.Request())
	return nil
}

func (h *Profile) Delete(ctx echo.Context) error {
	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)

	if err := h.auth.Logout(ctx); err != nil {
		log.Ctx(ctx).Error("error during logout on delete", "error", err)
	}

	if err := h.orm.User.DeleteOne(usr).Exec(ctx.Request().Context()); err != nil {
		return fail(err, "unable to delete user account")
	}

	msg.Success(ctx, "Your account has been deleted.")
	return redirect.New(ctx).Route(routenames.Home).Go()
}

func (h *Profile) AppearancePage(ctx echo.Context) error {
	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Settings/Appearance",
	)
}

func (h *Profile) PasswordPage(ctx echo.Context) error {
	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Settings/Password",
	)
}

func (h *Profile) UpdatePassword(ctx echo.Context) error {
	usr, ok := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	if !ok {
		msg.Danger(ctx, "You must be logged in.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	currentPasswordInput := ctx.FormValue("current_password")

	err := h.auth.CheckPassword(currentPasswordInput, usr.Password)
	if err != nil {
		msg.Danger(ctx, "The current password you entered is incorrect.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	password := ctx.FormValue("password")
	confirmPassword := ctx.FormValue("password_confirmation")

	err = h.auth.CheckPassword(password, confirmPassword)
	if err != nil {
		msg.Danger(ctx, "Password confirmation does not match. Please try again.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	_, err = h.orm.User.
		UpdateOneID(usr.ID).
		SetPassword(password).
		Save(ctx.Request().Context())
	if err != nil {
		msg.Danger(ctx, "Something went wrong while saving your new password.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	usr, err = h.orm.User.Get(ctx.Request().Context(), usr.ID)
	if err != nil {
		msg.Danger(ctx, "Something went wrong while refreshing your session.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	uri := ctx.Echo().Reverse(routenames.ProfileUpdatePassword)

	msg.Success(ctx, "Your password has been updated successfully.")
	h.Inertia.Redirect(ctx.Response().Writer, ctx.Request(), uri)
	return nil
}
