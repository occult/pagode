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
	profile.PATCH("/update", h.Update).Name = routenames.ProfileUpdate
	profile.DELETE("/delete", h.Delete).Name = routenames.ProfileDestroy

	profile.GET("/appearance", h.AppearancePage).Name = routenames.ProfileAppearance
	profile.GET("/password", h.PasswordPage).Name = routenames.ProfilePassword
}

func (h *Profile) EditPage(ctx echo.Context) error {
	usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)

	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Settings/Profile",
		inertia.Props{
			"user": map[string]any{
				"name":  usr.Name,
				"email": usr.Email,
			},
		},
	)
}

// TODO: Implement this later
func (h *Profile) Update(ctx echo.Context) error {
	// usr := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	//
	// var input forms.ProfileUpdate
	// if err := form.Submit(ctx, &input); err != nil {
	// 	switch err.(type) {
	// 	case validator.ValidationErrors:
	// 		return h.EditPage(ctx)
	// 	default:
	// 		return err
	// 	}
	// }
	//
	// _, err := usr.Update().
	// 	SetName(input.Name).
	// 	SetEmail(strings.ToLower(input.Email)).
	// 	Save(ctx.Request().Context())
	// if err != nil {
	// 	return fail(err, "unable to update user profile")
	// }

	msg.Success(ctx, "Your profile has been updated.")
	return redirect.New(ctx).Route(routenames.ProfileEdit).Go()
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
