package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/backlite/ui"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/admin"
	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/pager"
	"github.com/occult/pagode/pkg/redirect"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	"github.com/occult/pagode/pkg/ui/pages"
	inertia "github.com/romsar/gonertia/v2"
)

type Admin struct {
	orm      *ent.Client
	graph    *gen.Graph
	admin    *admin.Handler
	backlite *ui.Handler
	Inertia  *inertia.Inertia
}

func init() {
	Register(new(Admin))
}

func (h *Admin) Init(c *services.Container) error {
	var err error

	h.orm = c.ORM
	h.Inertia = c.Inertia
	h.graph = c.Graph
	h.orm = c.ORM
	h.admin = admin.NewHandler(h.orm, admin.HandlerConfig{
		ItemsPerPage: 25,
		PageQueryKey: pager.QueryKey,
		TimeFormat:   time.DateTime,
	})
	h.backlite, err = ui.NewHandler(ui.Config{
		DB:           c.Database,
		BasePath:     "/admin/tasks",
		ItemsPerPage: 25,
		ReleaseAfter: c.Config.Tasks.ReleaseAfter,
	})
	return err
}

func (h *Admin) Routes(g *echo.Group) {
	ag := g.Group("/admin/users")

	ag.GET("", h.Page).Name = routenames.AdminDashboard
	ag.POST("/add", h.AddUser).Name = routenames.AdminUserAdd
	ag.POST("/:id/edit", h.EditUser).Name = routenames.AdminUserEdit
	ag.POST("/:id/delete", h.DeleteUser).Name = routenames.AdminUserDelete

	entities := ag.Group("/entity")
	for _, n := range h.graph.Nodes {
		ng := entities.Group(fmt.Sprintf("/%s", strings.ToLower(n.Name)))
		ng.GET("", h.EntityList(n)).
			Name = routenames.AdminEntityList(n.Name)
		ng.GET("/add", h.EntityAdd(n)).
			Name = routenames.AdminEntityAdd(n.Name)
		ng.POST("/add", h.EntityAddSubmit(n)).
			Name = routenames.AdminEntityAddSubmit(n.Name)
		ng.GET("/:id/edit", h.EntityEdit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEdit(n.Name)
		ng.POST("/:id/edit", h.EntityEditSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEditSubmit(n.Name)
		ng.GET("/:id/delete", h.EntityDelete(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDelete(n.Name)
		ng.POST("/:id/delete", h.EntityDeleteSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDeleteSubmit(n.Name)
	}

	tasks := ag.Group("/tasks")
	tasks.GET("", h.Backlite(h.backlite.Running)).Name = routenames.AdminTasks
	tasks.GET("/succeeded", h.Backlite(h.backlite.Succeeded))
	tasks.GET("/failed", h.Backlite(h.backlite.Failed))
	tasks.GET("/upcoming", h.Backlite(h.backlite.Upcoming))
	tasks.GET("/task/:id", h.Backlite(h.backlite.Task))
	tasks.GET("/completed/:id", h.Backlite(h.backlite.TaskCompleted))
}

func (h *Admin) Page(ctx echo.Context) error {
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

func (h *Admin) AddUser(ctx echo.Context) error {
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

func (h *Admin) EditUser(ctx echo.Context) error {
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

func (h *Admin) DeleteUser(ctx echo.Context) error {
	w := ctx.Response().Writer
	r := ctx.Request()
	uri := ctx.Echo().Reverse("admin_dashboard")

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		msg.Danger(ctx, "Invalid user ID.")
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	err = h.orm.User.DeleteOneID(id).Exec(r.Context())
	if err != nil {
		msg.Danger(ctx, "Failed to delete user: "+err.Error())
		h.Inertia.Redirect(w, r, uri)
		return nil
	}

	msg.Success(ctx, "User successfully deleted.")
	h.Inertia.Redirect(w, r, uri)
	return nil
}

// middlewareEntityLoad is middleware to extract the entity ID and attempt to load the given entity.
func (h *Admin) middlewareEntityLoad(n *gen.Type) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid entity ID")
			}

			entity, err := h.admin.Get(ctx, n.Name, id)
			switch {
			case err == nil:
				ctx.Set(context.AdminEntityIDKey, id)
				ctx.Set(context.AdminEntityKey, map[string][]string(entity))
				return next(ctx)
			case ent.IsNotFound(err):
				return echo.NewHTTPError(http.StatusNotFound, "entity not found")
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}
}

func (h *Admin) EntityList(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		list, err := h.admin.List(ctx, n.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return pages.AdminEntityList(ctx, n.Name, list)
	}
}

func (h *Admin) EntityAdd(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), nil)
	}
}

func (h *Admin) EntityAddSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := h.admin.Create(ctx, n.Name)
		if err != nil {
			msg.Danger(ctx, err.Error())
			return h.EntityAdd(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Successfully added %s.", n.Name))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) EntityEdit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		v := ctx.Get(context.AdminEntityKey).(map[string][]string)
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), v)
	}
}

func (h *Admin) EntityEditSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Get(context.AdminEntityIDKey).(int)
		err := h.admin.Update(ctx, n.Name, id)
		if err != nil {
			msg.Danger(ctx, err.Error())
			return h.EntityEdit(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Updated %s.", n.Name))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) EntityDelete(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityDelete(ctx, n.Name)
	}
}

func (h *Admin) EntityDeleteSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Get(context.AdminEntityIDKey).(int)
		if err := h.admin.Delete(ctx, n.Name, id); err != nil {
			msg.Danger(ctx, err.Error())
			return h.EntityDelete(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Successfully deleted %s (ID %d).", n.Name, id))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) getEntitySchema(n *gen.Type) *load.Schema {
	for _, s := range h.graph.Schemas {
		if s.Name == n.Name {
			return s
		}
	}
	return nil
}

func (h *Admin) Backlite(handler func(http.ResponseWriter, *http.Request) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if id := c.Param("id"); id != "" {
			c.Request().SetPathValue("task", id)
		}
		return handler(c.Response().Writer, c.Request())
	}
}
