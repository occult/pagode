package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/session"
)

// Session sets the session storage in the request context
func Session(store sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			session.Store(ctx, store)
			return next(ctx)
		}
	}
}
