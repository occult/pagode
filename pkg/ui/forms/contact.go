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

type Contact struct {
	Email      string `form:"email" validate:"required,email"`
	Department string `form:"department" validate:"required,oneof=sales marketing hr"`
	Message    string `form:"message" validate:"required"`
	form.Submission
}

func (f *Contact) Render(r *ui.Request) Node {
	return Form(
		ID("contact"),
		Method(http.MethodPost),
		Attr("hx-post", r.Path(routenames.ContactSubmit)),
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Email",
			Name:      "email",
			InputType: "email",
			Label:     "Email address",
			Value:     f.Email,
		}),
		Radios(OptionsParams{
			Form:      f,
			FormField: "Department",
			Name:      "department",
			Label:     "Department",
			Value:     f.Department,
			Options: []Choice{
				{Value: "sales", Label: "Sales"},
				{Value: "marketing", Label: "Marketing"},
				{Value: "hr", Label: "HR"},
			},
		}),
		TextareaField(TextareaFieldParams{
			Form:      f,
			FormField: "Message",
			Name:      "message",
			Label:     "Message",
			Value:     f.Message,
		}),
		ControlGroup(
			FormButton("is-link", "Submit"),
		),
		CSRF(r),
	)
}
