package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"offer/controllers"
	"offer/models"
	"os"
	"pangpanglabs/goutils/config"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	appEnv = flag.String("app-env", os.Getenv("APP_ENV"), "app env")
)

func main() {
	var c struct {
		Database struct{ Driver, Connection string }
		Debug    bool
	}
	if err := config.Read(*appEnv, &c); err != nil {
		panic(err)
	}

	xormEngine, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	defer xormEngine.Close()
	xormEngine.ShowSQL(c.Debug)
	xormEngine.Sync(new(models.Discount))

	e := echo.New()
	e.Debug = c.Debug
	e.Renderer = NewTemplate()
	e.Validator = &Validator{}
	e.Static("/static", "static")

	controllers.HomeController{}.Init(e.Group("/"))
	controllers.DiscountController{}.Init(e.Group("/discounts"))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := xormEngine.NewSession()
			defer session.Close()
			c.Set("DB", &models.DB{session})

			if err := next(c); err != nil {
				session.Rollback()
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if c.Response().Status >= 500 {
				session.Rollback()
				return nil
			}
			err = session.Commit()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
	})
	e.Start(":8080")
}

type Validator struct {
}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

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
