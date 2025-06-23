package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

var handlers []Handler

// Handler handles one or more HTTP routes
type Handler interface {
	// Routes allows for self-registration of HTTP routes on the router
	Routes(g *echo.Group)

	// Init provides the service container to initialize
	Init(*services.Container) error
}

// Register registers a handler
func Register(h Handler) {
	handlers = append(handlers, h)
}

// GetHandlers returns all handlers
func GetHandlers() []Handler {
	return handlers
}

// fail is a helper to fail a request by returning a 500 error and logging the error
func fail(err error, log string, inertia *inertia.Inertia, c echo.Context) error {
	msg.Danger(c, fmt.Sprintf("%s: %v", log, err))

	req := c.Request()
	res := c.Response()

	inertia.Back(res, req)
	return nil
}
