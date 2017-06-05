package filters

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: func() *template.Template {
			templ := template.New("")
			if err := filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
				if strings.Contains(path, ".html") {
					_, err = templ.ParseFiles(path)
					if err != nil {
						log.Println(err)
					}
				}
				return err
			}); err != nil {
				panic(err)
			}
			return templ
		}(),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
