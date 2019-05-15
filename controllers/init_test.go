package controllers

import (
	"gopkg.in/testfixtures.v2"
	"runtime"

	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pangpanglabs/echosample/models"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	runtime.GOMAXPROCS(1)
	xormEngine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	xormEngine.Sync(new(models.Discount))

	fixtures, err := testfixtures.NewFolder(xormEngine.DB().DB, &testfixtures.SQLite{}, "../testdata/db_fixtures")
	if err != nil {
		panic(err)
	}
	testfixtures.SkipDatabaseNameCheck(true)

	if err := fixtures.Load(); err != nil {
		panic(err)
	}

	echoApp = echo.New()
	echoApp.Validator = &Validator{}

	logger := echomiddleware.ContextLogger()
	db := echomiddleware.ContextDB("test", xormEngine, echomiddleware.KafkaConfig{})

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return logger(db(handlerFunc))(c)
	}
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
