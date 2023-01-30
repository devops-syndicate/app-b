package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const APP_NAME = "app-b"
var jaeger_endpoint = os.Getenv("OTEL_EXPORTER_JAEGER_HTTP_ENDPOINT")

func main() {
  initLogger()
  tp, err := initTracer(jaeger_endpoint)
  if err != nil {
		logrus.Fatal(err)
	}

  otel.SetTracerProvider(tp)

  ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

  // Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logrus.Fatal(err)
		}
	}(ctx)
  
  gin.SetMode(gin.ReleaseMode)

  router := gin.New()

  router.Use(otelgin.Middleware(APP_NAME))

  // Setup route group for the API
  router.GET("/", HelloHandler)

  // Start and run the server
  router.Run(":8080")
}

func initLogger() {
  logrus.SetFormatter(&logrus.JSONFormatter{
    FieldMap: logrus.FieldMap{
      logrus.FieldKeyTime:  "timestamp",
      logrus.FieldKeyMsg:   "message",
  }})

  logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
    logrus.PanicLevel,
    logrus.FatalLevel,
    logrus.ErrorLevel,
    logrus.WarnLevel,
  )))
}

func initTracer(url string) (*sdktrace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
    sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(APP_NAME),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func HelloHandler(c *gin.Context) {
  logrus.WithContext(c.Request.Context()).Info("Callo hello endpoint")
 	c.String(http.StatusOK, "Hello '%s'", APP_NAME)
}