package forms

import (
	"net/http"

	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/ui"
	. "github.com/occult/pagode/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Cache struct {
	CurrentValue string
	Value        string `form:"value"`
	form.Submission
}

func (f *Cache) Render(r *ui.Request) Node {
	return Form(
		ID("cache"),
		Method(http.MethodPost),
		Attr("hx-post", r.Path(routenames.CacheSubmit)),
		Message(
			"is-info",
			"Test the cache",
			Group{
				P(Text("This route handler shows how the default in-memory cache works. Try updating the value using the form below and see how it persists after you reload the page.")),
				P(Text("HTMX makes it easy to re-render the cached value after the form is submitted.")),
			},
		),
		Label(
			For("value"),
			Class("value"),
			Text("Value in cache: "),
		),
		If(f.CurrentValue != "", Span(Class("tag is-success"), Text(f.CurrentValue))),
		If(f.CurrentValue == "", I(Text("(empty)"))),
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Value",
			Name:      "value",
			InputType: "text",
			Label:     "Value",
			Value:     f.Value,
		}),
		ControlGroup(
			FormButton("is-link", "Update cache"),
		),
		CSRF(r),
	)
}
