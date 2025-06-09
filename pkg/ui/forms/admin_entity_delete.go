package forms

import (
	"net/http"

	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/ui"
	. "github.com/occult/pagode/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminEntityDelete(r *ui.Request, entityTypeName string) Node {
	return Form(
		Method(http.MethodPost),
		P(
			Class("subtitle"),
			Textf("Are you sure you want to delete this %s?", entityTypeName),
		),
		ControlGroup(
			FormButton("is-link", "Delete"),
			ButtonLink(
				r.Path(routenames.AdminEntityList(entityTypeName)),
				"is-secondary",
				"Cancel",
			),
		),
		CSRF(r),
	)
}
