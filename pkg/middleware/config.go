package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/pkg/context"
)

// Config stores the configuration in the request so it can be accessed by the ui.
func Config(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(context.ConfigKey, cfg)
			return next(ctx)
		}
	}
}
