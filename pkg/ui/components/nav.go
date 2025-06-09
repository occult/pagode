package components

import (
	"github.com/occult/pagode/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MenuLink(r *ui.Request, title, routeName string, routeParams ...any) Node {
	href := r.Path(routeName, routeParams...)

	return Li(
		A(
			Href(href),
			Text(title),
			If(href == r.CurrentPath, Class("is-active")),
		),
	)
}
