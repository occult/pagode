package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/pkg/log"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/redirect"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"

	inertia "github.com/romsar/gonertia/v2"
)

type Auth struct {
	config  *config.Config
	auth    *services.AuthClient
	casdoor *services.CasdoorClient
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Auth))
}

func (h *Auth) Init(c *services.Container) error {
	h.config = c.Config
	h.auth = c.Auth
	h.casdoor = c.Casdoor
	h.Inertia = c.Inertia
	return nil
}

func (h *Auth) Routes(g *echo.Group) {
	g.GET("/logout", h.Logout, middleware.RequireAuthentication).Name = routenames.Logout

	noAuth := g.Group("/user", middleware.RequireNoAuthentication)
	noAuth.GET("/login", h.LoginPage).Name = routenames.Login
	noAuth.GET("/register", h.RegisterPage).Name = routenames.Register
}

func (h *Auth) LoginPage(ctx echo.Context) error {
	if !h.casdoor.IsReachable() {
		log.Ctx(ctx).Error("authentication service unreachable")
		msg.Danger(ctx, "Authentication is temporarily unavailable. Please try again later.")
		return redirect.New(ctx).Route(routenames.Welcome).Go()
	}

	callbackURL := h.config.App.Host + ctx.Echo().Reverse(routenames.CasdoorCallback)
	signinURL := h.casdoor.GetSigninURL(callbackURL, "pagode")

	if ctx.Request().Header.Get("X-Inertia") != "" {
		ctx.Response().Header().Set("X-Inertia-Location", signinURL)
		return ctx.NoContent(http.StatusConflict)
	}
	return ctx.Redirect(http.StatusSeeOther, signinURL)
}

func (h *Auth) RegisterPage(ctx echo.Context) error {
	if !h.casdoor.IsReachable() {
		log.Ctx(ctx).Error("authentication service unreachable")
		msg.Danger(ctx, "Authentication is temporarily unavailable. Please try again later.")
		return redirect.New(ctx).Route(routenames.Welcome).Go()
	}

	callbackURL := h.config.App.Host + ctx.Echo().Reverse(routenames.CasdoorCallback)
	signupURL := h.casdoor.GetSignupURL(callbackURL, "pagode")

	if ctx.Request().Header.Get("X-Inertia") != "" {
		ctx.Response().Header().Set("X-Inertia-Location", signupURL)
		return ctx.NoContent(http.StatusConflict)
	}
	return ctx.Redirect(http.StatusSeeOther, signupURL)
}

func (h *Auth) Logout(ctx echo.Context) error {
	if err := h.auth.Logout(ctx); err == nil {
		msg.Success(ctx, "You have been logged out successfully.")
	} else {
		msg.Danger(ctx, "An error occurred. Please try again.")
	}

	// Redirect to Casdoor logout to clear the SSO session as well.
	if h.casdoor.IsReachable() {
		logoutURL := h.casdoor.GetLogoutURL()
		if logoutURL != "" {
			return ctx.Redirect(http.StatusSeeOther, logoutURL)
		}
	}

	return redirect.New(ctx).
		Route(routenames.Welcome).
		Go()
}
