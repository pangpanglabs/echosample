package filters

import (
	"context"
	"log"
	"net/http"
	"echosample/factory"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
)

func SetDbContext(db *xorm.Engine) echo.MiddlewareFunc {
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
