package handlers

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	assets "github.com/occult/pagode"
	"github.com/occult/pagode/pkg/services"
	"github.com/spf13/afero"
)

type Build struct {
	files afero.Fs
}

func init() {
	Register(new(Build))
}

func (h *Build) Init(c *services.Container) error {
	return nil
}

func (h *Build) Routes(g *echo.Group) {
	sub, err := fs.Sub(assets.StaticFS, "public/build/assets")
	if err != nil {
		panic(err)
	}

	handler := http.StripPrefix("/build/assets/", http.FileServer(http.FS(sub)))
	g.GET("/build/assets/*", echo.WrapHandler(handler))
}
