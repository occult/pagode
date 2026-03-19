package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/user"
	"github.com/occult/pagode/pkg/log"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/redirect"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Casdoor struct {
	auth    *services.AuthClient
	casdoor *services.CasdoorClient
	orm     *ent.Client
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Casdoor))
}

func (h *Casdoor) Init(c *services.Container) error {
	h.auth = c.Auth
	h.casdoor = c.Casdoor
	h.orm = c.ORM
	h.Inertia = c.Inertia
	return nil
}

func (h *Casdoor) Routes(g *echo.Group) {
	if h.casdoor == nil {
		return
	}
	g.GET("/auth/casdoor/callback", h.Callback).Name = routenames.CasdoorCallback
}

func (h *Casdoor) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	state := ctx.QueryParam("state")

	if code == "" {
		msg.Danger(ctx, "Invalid authentication response.")
		return redirect.New(ctx).Route(routenames.Welcome).Go()
	}

	email, name, err := h.casdoor.ExchangeCodeAndGetUser(code, state)
	if err != nil {
		log.Ctx(ctx).Error("casdoor auth failed", "error", err)
		msg.Danger(ctx, "Authentication failed. Please try again.")
		return redirect.New(ctx).Route(routenames.Welcome).Go()
	}

	// Find or create local user
	u, err := h.orm.User.
		Query().
		Where(user.Email(strings.ToLower(email))).
		Only(ctx.Request().Context())

	if err != nil {
		if ent.IsNotFound(err) {
			// Create a new local user with a random placeholder password
			placeholder, _ := randomPassword(64)
			u, err = h.orm.User.
				Create().
				SetName(name).
				SetEmail(strings.ToLower(email)).
				SetPassword(placeholder).
				SetVerified(true).
				Save(ctx.Request().Context())
			if err != nil {
				log.Ctx(ctx).Error("failed to create user from casdoor", "error", err)
				msg.Danger(ctx, "Failed to create account. Please try again.")
				return redirect.New(ctx).Route(routenames.Welcome).Go()
			}
		} else {
			log.Ctx(ctx).Error("failed to query user", "error", err)
			msg.Danger(ctx, "An error occurred. Please try again.")
			return redirect.New(ctx).Route(routenames.Welcome).Go()
		}
	}

	// Mark as verified if not already
	if !u.Verified {
		u, _ = u.Update().SetVerified(true).Save(ctx.Request().Context())
	}

	// Log the user in using the existing session mechanism
	if err := h.auth.Login(ctx, u.ID); err != nil {
		log.Ctx(ctx).Error("failed to login casdoor user", "error", err)
		msg.Danger(ctx, "Login failed. Please try again.")
		return redirect.New(ctx).Route(routenames.Welcome).Go()
	}

	msg.Success(ctx, "You have been logged in successfully.")
	return redirect.New(ctx).Route(routenames.Dashboard).Go()
}

func randomPassword(length int) (string, error) {
	b := make([]byte, (length/2)+1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:length], nil
}
