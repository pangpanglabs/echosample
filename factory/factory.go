package factory

import (
	"context"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

func DB(ctx context.Context) *xorm.Session {
	v := ctx.Value(echomiddleware.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*xorm.Session); ok {
		return db
	}
	if db, ok := v.(*xorm.Engine); ok {
		return db.NewSession()
	}
	panic("DB is not exist")
}

func BehaviorLogger(ctx context.Context) *echomiddleware.BehaviorLogContext {
	v := ctx.Value(echomiddleware.BehaviorLoggerName)
	if logger, ok := v.(*echomiddleware.BehaviorLogContext); ok {
		return logger
	}
	return echomiddleware.NewNopLogger()
}

func Logger(ctx context.Context) *logrus.Entry {
	v := ctx.Value(echomiddleware.ContextLoggerName)
	if v == nil {
		return logrus.WithFields(logrus.Fields{})
	}
	if logger, ok := v.(*logrus.Entry); ok {
		return logger
	}
	return logrus.WithFields(logrus.Fields{})
}
