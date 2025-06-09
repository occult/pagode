package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
	"github.com/spf13/afero"
)

type Files struct {
	files   afero.Fs
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Files))
}

func (h *Files) Init(c *services.Container) error {
	h.files = c.Files
	h.Inertia = c.Inertia
	return nil
}

func (h *Files) Routes(g *echo.Group) {
	g.GET("/files", h.UploadFilePage).Name = routenames.Files
	g.POST("/files", h.Submit).Name = routenames.FilesSubmit
}

func (h *Files) UploadFilePage(ctx echo.Context) error {
	info, err := afero.ReadDir(h.files, "")
	if err != nil {
		return err
	}

	files := make([]map[string]interface{}, 0, len(info))
	for _, file := range info {
		files = append(files, map[string]interface{}{
			"id":       file.Name(),
			"name":     file.Name(),
			"size":     file.Size(),
			"modified": file.ModTime().Format(time.DateTime),
		})
	}

	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"UploadFile",
		inertia.Props{
			"files": files,
		},
	)
}

func (h *Files) Submit(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		msg.Danger(ctx, "A file is required.")
		return h.UploadFilePage(ctx)
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := h.files.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	msg.Success(ctx, fmt.Sprintf("%s was uploaded successfully.", file.Filename))

	h.Inertia.Redirect(
		ctx.Response().Writer,
		ctx.Request(),
		ctx.Echo().Reverse(routenames.Files),
	)

	return nil
}
