package handlers

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	"github.com/occult/pagode/pkg/ui/forms"
	"github.com/occult/pagode/pkg/ui/pages"
	inertia "github.com/romsar/gonertia/v2"
)

type Cache struct {
	cache   *services.CacheClient
	inertia *inertia.Inertia
}

func init() {
	Register(new(Cache))
}

func (h *Cache) Init(c *services.Container) error {
	h.cache = c.Cache
	h.inertia = c.Inertia
	return nil
}

func (h *Cache) Routes(g *echo.Group) {
	g.GET("/cache", h.Page).Name = routenames.Cache
	g.POST("/cache", h.Submit).Name = routenames.CacheSubmit
}

func (h *Cache) Page(ctx echo.Context) error {
	f := form.Get[forms.Cache](ctx)

	// Fetch the value from the cache.
	value, err := h.cache.
		Get().
		Key("page_cache_example").
		Fetch(ctx.Request().Context())

	// Store the value in the form, so it can be rendered, if found.
	switch {
	case err == nil:
		f.CurrentValue = value.(string)
	case errors.Is(err, services.ErrCacheMiss):
	default:
		return fail(err, "failed to fetch from cache", h.inertia, ctx)
	}

	return pages.UpdateCache(ctx, f)
}

func (h *Cache) Submit(ctx echo.Context) error {
	var input forms.Cache

	if err := form.Submit(ctx, &input); err != nil {
		return err
	}

	// Set the cache.
	err := h.cache.
		Set().
		Key("page_cache_example").
		Data(input.Value).
		Expiration(30 * time.Minute).
		Save(ctx.Request().Context())
	if err != nil {
		return fail(err, "failed to fetch from cache", h.inertia, ctx)
	}

	form.Clear(ctx)

	return h.Page(ctx)
}
