package main

import (
	"fmt"
	"go_uploadfile/models"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Instantiate a template registry and register all html files inside the view folder
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.POST("/upload", func(c echo.Context) error {
		err := os.Mkdir("./static/images", 0755)

		if err != nil {
			fmt.Println("Cannot create a file when that file already exists.")
		}

		file, err := c.FormFile("uploadfile")
		if err != nil {
			return c.JSON(http.StatusOK, models.Response{Code: 400, Message: "file not have data"})
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create("./static/images/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, models.Response{Code: 200, Message: "upload  successfully !"})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
