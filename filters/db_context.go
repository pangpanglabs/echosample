package filters

import (
	"net/http"
	"offer/models"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
)

func SetDbContext(db *xorm.Engine) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := db.NewSession()
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
			if err := session.Commit(); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
	}
}
