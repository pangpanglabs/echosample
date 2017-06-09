package filters

import (
	"context"
	"log"
	"net/http"
	"offer/models"
	"runtime"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"

	"github.com/pangpanglabs/echosample/factory"
)

func SetDbContext(db *xorm.Engine) echo.MiddlewareFunc {
	if db.Dialect().DriverName() == "sqlite3" {
		// sqlite does not support concurrency
		runtime.GOMAXPROCS(1)
	}
	db.ShowSQL(true)
	db.Sync(new(models.Discount))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := db.NewSession()
			defer session.Close()

			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), factory.ContextDBName, session)))

			switch req.Method {
			case "POST", "PUT", "DELETE":
				if err := session.Begin(); err != nil {
					log.Println(err)
				}
				if err := next(c); err != nil {
					session.Rollback()
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				if err := next(c); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}

			return nil
		}
	}
}
