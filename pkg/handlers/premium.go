package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Premium struct {
	Inertia *inertia.Inertia
	Auth    *services.AuthClient
	ORM     *ent.Client
}

func init() {
	Register(new(Premium))
}

func (h *Premium) Init(c *services.Container) error {
	h.Inertia = c.Inertia
	h.Auth = c.Auth
	h.ORM = c.ORM
	return nil
}

func (h *Premium) Routes(g *echo.Group) {
	authGroup := g.Group("")
	authGroup.Use(middleware.RequireAuthentication)
	authGroup.Use(middleware.RequirePaidUser(h.ORM))
	
	authGroup.GET("/premium", h.Page).Name = routenames.Premium
}

func (h *Premium) Page(ctx echo.Context) error {
	user, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Premium",
		inertia.Props{
			"title": "Premium Access",
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	)
}