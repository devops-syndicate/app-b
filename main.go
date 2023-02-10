package main

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

const APP_NAME = "app-b"

func main() {
	initLogger()
	initProfiling()
	tp := initTracing()
	defer shutdownTracing(tp)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	p := ginprometheus.NewPrometheus("http")
	p.Use(router)

	router.Use(otelgin.Middleware(APP_NAME))

	router.GET("/", HelloHandler)
	router.GET("/random", randomHandler(tp.Tracer("")))

	logrus.Info("Start listening on port 8080")

	// Start and run the server
	router.Run(":8080")
}

func HelloHandler(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("Hello endpoint called")
	c.String(http.StatusOK, "Hello '%s'", APP_NAME)
}

func randomHandler(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithContext(c.Request.Context()).Info("Random endpoint called")
		callHttpbin(c.Request.Context(), tracer)
		c.String(http.StatusOK, "From random")
	}
}

func callHttpbin(c context.Context, tracer trace.Tracer) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	logrus.WithContext(c).Infof("Call httpbin with %d seconds delay", n)
	ctx, span := tracer.Start(c, "call httpbin")
	defer span.End()
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	resp, _ := otelhttp.Get(ctx, "http://httpbin.httpbin/delay/"+strconv.Itoa(n))
	defer resp.Body.Close()
}
