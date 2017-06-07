package filters

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"offer/factory"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func SetLogger(env string) echo.MiddlewareFunc {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logger.Writer())

	if env == "production" {
		logger.Formatter = &logrus.JSONFormatter{}
		logger.Level = logrus.InfoLevel

		hooks := logrus.LevelHooks{}
		hooks.Add(&CallkerHook{})
		logger.Hooks = hooks
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			if env == "production" {
				id := c.Request().Header.Get(echo.HeaderXRequestID)
				if id == "" {
					id = c.Response().Header().Get(echo.HeaderXRequestID)
				}
				logEntry = logEntry.WithField("request_id", id)
			}

			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), factory.ContextLoggerName, logEntry)))

			return next(c)
		}
	}
}

type CallkerHook struct{}

func (c *CallkerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}
func (c *CallkerHook) Fire(entry *logrus.Entry) error {
	var ok bool
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}
	entry.Data["caller"] = fmt.Sprintf("%s:%d", file, line)
	return nil
}
