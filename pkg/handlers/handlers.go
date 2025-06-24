package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/services"
)

var handlers []Handler

// Handler handles one or more HTTP routes
type Handler interface {
	// Routes allows for self-registration of HTTP routes on the router
	Routes(g *echo.Group)

	// Init provides the service container to initialize
	Init(*services.Container) error
}

// InertiaBacker abstracts the Back method from gonertia.Inertia
// to allow injection and mocking in handlers and tests.
type InertiaBacker interface {
	Back(http.ResponseWriter, *http.Request, ...int)
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
func fail(err error, log string, inertia InertiaBacker, c echo.Context) error {
	msg.Danger(c, fmt.Sprintf("%s: %v", log, err))

	req := c.Request()
	res := c.Response()

	inertia.Back(res, req)
	return nil
}
