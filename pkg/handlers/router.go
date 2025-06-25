package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/services"
)

// BuildRouter builds the router.
func BuildRouter(c *services.Container) error {
	// Static files with proper cache control
	staticGroup := c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile))
	
	// Standard static files serving
	staticGroup.Static(config.StaticPrefix, config.StaticDir)
	
	// Assets serving - unified path for all environments
	if c.Config.App.Environment == config.EnvProduction {
		staticGroup.Static("/files", "/app/static")
	} else {
		staticGroup.Static("/files", filepath.Join(services.ProjectRoot(), "public"))
	}

	// Non-static file route group.
	g := c.Web.Group("")

	// Force HTTPS, if enabled.
	if c.Config.HTTP.TLS.Enabled {
		g.Use(echomw.HTTPSRedirect())
	}

	// Create a cookie store for session data.
	cookieStore := sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.SameSite = http.SameSiteStrictMode

	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		middleware.SetLogger(),
		middleware.LogRequest(),
		echomw.Gzip(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
		middleware.Config(c.Config),
		middleware.Session(cookieStore),
		middleware.LoadAuthenticatedUser(c.Auth),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup:    "header:X-XSRF-TOKEN", // where to look for token
			CookieName:     "XSRF-TOKEN",          // this sets the cookie
			CookiePath:     "/",                   // make it accessible app-wide
			CookieHTTPOnly: false,                 // must be false so JS (Axios) can read it
			CookieSameSite: http.SameSiteStrictMode,
			ContextKey:     context.CSRFKey,
		}),
		middleware.InertiaProps(), // leave this as the last one
	)

	// Error handler.
	errHandler := &Error{}
	_ = errHandler.Init(c)

	c.Web.HTTPErrorHandler = func(err error, ctx echo.Context) {
		errHandler.Page(err, ctx)
	}

	// Initialize and register all handlers.
	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}

		h.Routes(g)
	}

	return nil
}
