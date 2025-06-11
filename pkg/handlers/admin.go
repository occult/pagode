package handlers

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type AdminUser struct {
	orm     *ent.Client
	Inertia *inertia.Inertia
}

func init() {
	Register(new(AdminUser))
}

func (h *AdminUser) Init(c *services.Container) error {
	h.orm = c.ORM
	h.Inertia = c.Inertia
	return nil
}

func (h *AdminUser) Routes(g *echo.Group) {
	ag := g.Group("/admin/users")
	ag.GET("", h.Page).Name = "admin_dashboard"
	ag.POST("/add", h.AddUser).Name = "admin_user_add"
	ag.POST("/:id/edit", h.EditUser).Name = "admin_user_edit"
}

func (h *AdminUser) Page(ctx echo.Context) error {
	pageStr := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	total, err := h.orm.User.Query().Count(ctx.Request().Context())
	if err != nil {
		msg.Danger(ctx, "Failed to count users.")
		return ctx.NoContent(500)
	}

	users, err := h.orm.User.
		Query().
		Limit(limit).
		Offset(offset).
		All(ctx.Request().Context())
	if err != nil {
		msg.Danger(ctx, "Failed to load users.")
		return ctx.NoContent(500)
	}

	totalPages := (total + limit - 1) / limit

	err = h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Admin/AdminView",
		inertia.Props{
			"users": users,
			"pagination": map[string]any{
				"total":      total,
				"page":       page,
				"perPage":    limit,
				"totalPages": totalPages,
			},
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

type AdminUserForm struct {
	Name          string `form:"name" validate:"required"`
	Email         string `form:"email" validate:"required,email"`
	Admin         bool   `form:"admin"`
	EmailVerified bool   `form:"emailVerified"`
	form.Submission
}

func (h *AdminUser) AddUser(ctx echo.Context) error {
	w := ctx.Response().Writer
	r := ctx.Request()
	uri := ctx.Echo().Reverse("admin_dashboard")

	var input AdminUserForm
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case validator.ValidationErrors:
		msg.Danger(ctx, "Please fill in all fields correctly.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	case nil:
	default:
		msg.Danger(ctx, "Invalid form data.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	_, err = h.orm.User.
		Create().
		SetName(input.Name).
		SetEmail(strings.ToLower(input.Email)).
		SetAdmin(input.Admin).
		SetVerified(input.EmailVerified).
		Save(r.Context())
	if err != nil {
		msg.Danger(ctx, "Failed to create user: "+err.Error())
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	msg.Success(ctx, "User successfully created.")
	h.Inertia.Redirect(w, r, uri)
	return nil
}

func (h *AdminUser) EditUser(ctx echo.Context) error {
	w := ctx.Response().Writer
	r := ctx.Request()
	uri := ctx.Echo().Reverse("admin_dashboard")

	var input AdminUserForm
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case validator.ValidationErrors:
		msg.Danger(ctx, "Please fill in all fields correctly.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	case nil:
	default:
		msg.Danger(ctx, "Invalid form data.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	id, convErr := strconv.Atoi(ctx.Param("id"))
	if convErr != nil {
		msg.Danger(ctx, "Invalid user ID.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	err = h.orm.User.
		UpdateOneID(id).
		SetName(input.Name).
		SetEmail(strings.ToLower(input.Email)).
		SetAdmin(input.Admin).
		SetVerified(input.EmailVerified).
		Exec(r.Context())
	if err != nil {
		msg.Danger(ctx, "Failed to update user: "+err.Error())
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	msg.Success(ctx, "User successfully updated.")
	h.Inertia.Redirect(w, r, uri)
	return nil
}
