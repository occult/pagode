package handlers

import (
	"fmt"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/ui/forms"
	"github.com/occult/pagode/pkg/ui/pages"
	inertia "github.com/romsar/gonertia/v2"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/services"
	"github.com/occult/pagode/pkg/tasks"
)

type Task struct {
	tasks   *backlite.Client
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Task))
}

func (h *Task) Init(c *services.Container) error {
	h.tasks = c.Tasks
	return nil
}

func (h *Task) Routes(g *echo.Group) {
	g.GET("/task", h.Page).Name = routenames.Task
	g.POST("/task", h.Submit).Name = routenames.TaskSubmit
}

func (h *Task) Page(ctx echo.Context) error {
	return pages.AddTask(ctx, form.Get[forms.Task](ctx))
}

func (h *Task) Submit(ctx echo.Context) error {
	var input forms.Task

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.Page(ctx)
	default:
		return err
	}

	// Insert the task
	err = h.tasks.
		Add(tasks.ExampleTask{
			Message: input.Message,
		}).
		Wait(time.Duration(input.Delay) * time.Second).
		Save()
	if err != nil {
		return fail(err, "unable to create a task", h.Inertia, ctx)
	}

	msg.Success(ctx, fmt.Sprintf("The task has been created. Check the logs in %d seconds.", input.Delay))
	form.Clear(ctx)

	return h.Page(ctx)
}
