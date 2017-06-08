package filters

import (
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
)

func Tracer(zipkinAddr, hostPort, serviceName string, debug bool) echo.MiddlewareFunc {
	if zipkinAddr == "" {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error { return next(c) }
		}
	}

	collector, err := zipkin.NewHTTPCollector(zipkinAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	// defer collector.Close()

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, debug, hostPort, serviceName),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.WithFields(logrus.Fields{
		"tracer": "ZipkinHTTP",
		"addr":   zipkinAddr,
	}).Info("Set Tracer")

	operationName := "http"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(req.Header),
			)
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				logrus.Error(err)
			}
			span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ext.HTTPMethod.Set(span, req.Method)
			ext.HTTPUrl.Set(span, req.URL.String())
			ext.SpanKindRPCServer.Set(span)

			c.SetRequest(req.WithContext(opentracing.ContextWithSpan(req.Context(), span)))

			return next(c)
		}
	}
}
