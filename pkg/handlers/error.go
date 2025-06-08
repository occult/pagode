package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/log"
	"github.com/mikestefanello/pagoda/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Error struct {
	Inertia *inertia.Inertia
}

func (e *Error) Init(c *services.Container) error {
	e.Inertia = c.Inertia
	return nil
}

func (e *Error) Page(err error, ctx echo.Context) {
	if ctx.Response().Committed || context.IsCanceledError(err) {
		return
	}

	// Determine status code
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// Log based on error type
	logger := log.Ctx(ctx)
	switch {
	case code >= 500:
		logger.Error(err.Error())
	case code >= 400:
		logger.Warn(err.Error())
	}

	// Write status code header
	ctx.Response().WriteHeader(code)

	// Render Inertia error page
	renderErr := e.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"ErrorPage",
		inertia.Props{
			"status": code,
		},
	)
	if renderErr != nil {
		log.Ctx(ctx).Error("failed to render error page", "error", renderErr)
	}
}
