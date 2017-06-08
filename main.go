package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/goutils/config"

	"github.com/pangpanglabs/echosample/controllers"
	"github.com/pangpanglabs/echosample/filters"
	"github.com/pangpanglabs/echosample/models"
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	flag.Parse()

	var c struct {
		Database struct{ Driver, Connection string }
		Debug    bool
		Httpport string
	}
	if err := config.Read(*appEnv, &c); err != nil {
		panic(err)
	}

	xormEngine, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	defer xormEngine.Close()
	if c.Database.Driver == "sqlite3" {
		// sqlite does not support concurrency
		runtime.GOMAXPROCS(1)
	}
	xormEngine.ShowSQL(c.Debug)
	xormEngine.Sync(new(models.Discount))

	e := echo.New()

	controllers.HomeController{}.Init(e.Group("/"))
	controllers.DiscountController{}.Init(e.Group("/discounts"))
	controllers.DiscountApiController{}.Init(e.Group("/api/discounts"))

	e.Static("/static", "static")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(filters.SetDbContext(xormEngine))
	e.Use(filters.SetLogger(*appEnv))

	e.Renderer = filters.NewTemplate()
	e.Validator = &filters.Validator{}
	e.Debug = c.Debug

	if err := e.Start(":" + c.Httpport); err != nil {
		log.Println(err)
	}

}
