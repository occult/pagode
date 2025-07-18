package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/ui"
	. "github.com/occult/pagode/pkg/ui/components"
	"github.com/occult/pagode/pkg/ui/forms"
	"github.com/occult/pagode/pkg/ui/layouts"
	"github.com/occult/pagode/pkg/ui/models"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UploadFile(ctx echo.Context, files []*models.File) error {
	r := ui.NewRequest(ctx)
	r.Title = "Upload a file"

	fileList := make(Group, len(files))
	for i, file := range files {
		fileList[i] = file.Render()
	}

	n := Group{
		Message(
			"is-link",
			"",
			P(Text("This is a very basic example of how to handle file uploads. Files uploaded will be saved to the directory specified in your configuration.")),
		),
		Hr(),
		forms.File{}.Render(r),
		Hr(),
		H3(
			Class("title"),
			Text("Uploaded files"),
		),
		Message("is-warning", "", P(Text("Below are all files in the configured upload directory."))),
		Table(
			Class("table"),
			THead(
				Tr(
					Th(Text("Filename")),
					Th(Text("Size")),
					Th(Text("Modified on")),
				),
			),
			TBody(
				fileList,
			),
		),
	}

	return r.Render(layouts.Primary, n)
}
