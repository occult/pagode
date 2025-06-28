package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/services"
	"github.com/spf13/afero"
)

type Build struct {
	files      afero.Fs
	buildAssets embed.FS
}

func init() {
	Register(new(Build))
}

func (h *Build) Init(c *services.Container) error {
	return nil
}

func (h *Build) SetBuildAssets(assets embed.FS) {
	h.buildAssets = assets
}

func (h *Build) Routes(g *echo.Group) {
	// Serve the embedded build directory
	distFS, err := fs.Sub(h.buildAssets, "dist")
	if err != nil {
		panic(err)
	}
	fs := http.StripPrefix("/build/", http.FileServer(http.FS(distFS)))
	g.GET("/build/*", echo.WrapHandler(fs))
}
