package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	"github.com/occult/pagode/pkg/ui/forms"
	"github.com/occult/pagode/pkg/ui/pages"
	inertia "github.com/romsar/gonertia/v2"
)

type Contact struct {
	mail    *services.MailClient
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Contact))
}

func (h *Contact) Init(c *services.Container) error {
	h.mail = c.Mail
	return nil
}

func (h *Contact) Routes(g *echo.Group) {
	g.GET("/contact", h.Page).Name = routenames.Contact
	g.POST("/contact", h.Submit).Name = routenames.ContactSubmit
}

func (h *Contact) Page(ctx echo.Context) error {
	return pages.ContactUs(ctx, form.Get[forms.Contact](ctx))
}

func (h *Contact) Submit(ctx echo.Context) error {
	var input forms.Contact

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.Page(ctx)
	default:
		return err
	}

	err = h.mail.
		Compose().
		To(input.Email).
		Subject("Contact form submitted").
		Body(fmt.Sprintf("The message is: %s", input.Message)).
		Send(ctx)
	if err != nil {
		return fail(err, "unable to send email", h.Inertia, ctx)
	}

	return h.Page(ctx)
}
