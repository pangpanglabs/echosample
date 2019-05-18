package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/echoswagger"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/echotpl"

	"github.com/pangpanglabs/echosample/controllers"
	"github.com/pangpanglabs/echosample/models"
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	flag.Parse()

	var c Config
	if err := configutil.Read(*appEnv, &c); err != nil {
		panic(err)
	}
	fmt.Println(c)
	db, err := initDB(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	r := echoswagger.New(e, "/doc", &echoswagger.Info{
		Title:       "Echo Sample",
		Description: "This is API doc for Echo Sample",
		Version:     "1.0",
	})

	controllers.HomeController{}.Init(e.Group("/"))
	controllers.DiscountController{}.Init(e.Group("/discounts"))
	controllers.DiscountApiController{}.Init(r.Group("Discount", "/api/discounts"))

	e.Static("/static", "static")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(echomiddleware.ContextBase())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// e.Use(middleware.Logger())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDB(c.Service, db, echomiddleware.KafkaConfig(c.Database.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(c.Service, echomiddleware.KafkaConfig(c.BehaviorLog.Kafka)))

	e.Renderer = echotpl.New()
	e.Validator = &Validator{}
	e.Debug = c.Debug

	if err := e.Start(":" + c.HttpPort); err != nil {
		log.Println(err)
	}

}

func initDB(driver, connection string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	db.Sync(new(models.Discount))
	return db, nil
}

type Config struct {
	Database struct {
		Driver     string
		Connection string
		Logger     struct {
			Kafka echomiddleware.KafkaConfig
		}
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	Trace struct {
		Zipkin echomiddleware.ZipkinConfig
	}

	Debug    bool
	Service  string
	HttpPort string
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
