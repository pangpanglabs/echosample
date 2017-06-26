package filters

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"

	"github.com/pangpanglabs/echosample/config"
	"github.com/pangpanglabs/echosample/factory"
	"github.com/pangpanglabs/echosample/models"
	"github.com/pangpanglabs/goutils/kafka"
)

func DbContext(c config.Database) echo.MiddlewareFunc {
	db, err := xorm.NewEngine(c.Driver, c.Connection)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	if db.Dialect().DriverName() == "sqlite3" {
		// sqlite does not support concurrency
		runtime.GOMAXPROCS(1)
	}

	if len(c.Logger.Kafka.Brokers) != 0 {
		if producer, err := kafka.NewProducer(c.Logger.Kafka.Brokers, c.Logger.Kafka.Topic, func(c *sarama.Config) {
			c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
			c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
			c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

		}); err == nil {
			db.SetLogger(&dbLogger{serviceName: config.Const.ServiceName, Producer: producer})
		}
	}

	db.ShowSQL()
	db.ShowExecTime()
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
					return err
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				return next(c)
			}

			return nil
		}
	}
}

type dbLogger struct {
	serviceName string
	*kafka.Producer
}

func (logger *dbLogger) Write(v []interface{}) {
	if len(v) == 3 {
		logger.Send(map[string]interface{}{
			"service": logger.serviceName,
			"sql":     v[0],
			"args":    v[1],
			"took":    v[2],
		})
	} else if len(v) == 2 {
		logger.Send(map[string]interface{}{
			"service": logger.serviceName,
			"sql":     v[0],
			"took":    v[1],
		})
	}
}
func (logger *dbLogger) Infof(format string, v ...interface{})  { logger.Write(v) }
func (logger *dbLogger) Errorf(format string, v ...interface{}) {}
func (logger *dbLogger) Debugf(format string, v ...interface{}) {}
func (logger *dbLogger) Warnf(format string, v ...interface{})  {}

func (logger *dbLogger) Debug(v ...interface{})   {}
func (logger *dbLogger) Error(v ...interface{})   {}
func (logger *dbLogger) Info(v ...interface{})    {}
func (logger *dbLogger) Warn(v ...interface{})    {}
func (logger *dbLogger) SetLevel(l core.LogLevel) {}
func (logger *dbLogger) ShowSQL(show ...bool)     {}
func (logger *dbLogger) Level() core.LogLevel     { return 0 }
func (logger *dbLogger) IsShowSQL() bool          { return true }
