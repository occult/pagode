package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/ui"
	"github.com/occult/pagode/pkg/ui/forms"
	"github.com/occult/pagode/pkg/ui/layouts"
)

func UpdateCache(ctx echo.Context, form *forms.Cache) error {
	r := ui.NewRequest(ctx)
	r.Title = "Set a cache entry"

	return r.Render(layouts.Primary, form.Render(r))
}
