package models_test

import (
	"context"
	"offer/factory"
	"offer/models"
	"runtime"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var ctx context.Context

func init() {
	runtime.GOMAXPROCS(1)
	xormEngine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	xormEngine.ShowSQL(true)
	xormEngine.Sync(new(models.Discount))
	ctx = context.WithValue(context.Background(), factory.ContextDBName, xormEngine.NewSession())
}
