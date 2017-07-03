package models

import (
	"context"
	"runtime"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pangpanglabs/echosample/factory"
)

var ctx context.Context

func init() {
	runtime.GOMAXPROCS(1)
	xormEngine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	xormEngine.ShowSQL(true)
	xormEngine.Sync(new(Discount))
	ctx = context.WithValue(context.Background(), factory.ContextDBName, xormEngine.NewSession())
}
