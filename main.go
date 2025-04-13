package main

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const APP_NAME = "app-b"

func main() {
	initLogger()
	initProfiling()
	tp := initTracing()
	defer shutdownTracing(tp)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			HandlerIDLabel: "route",
		}),
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(otelgin.Middleware(APP_NAME))
	router.Use(ginmiddleware.Handler("", mdlw))

	router.GET("/", HelloHandler)
	router.GET("/random", RandomHandler)

	logrus.Info("Start listening on port 8080")

	// Start and run the server
	router.Run(":8080")
}

func HelloHandler(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("Hello endpoint called")
	c.String(http.StatusOK, "Hello '%s'", APP_NAME)
}

func RandomHandler(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("Random endpoint called")
	callHttpbin(c.Request.Context())
	c.String(http.StatusOK, "From random")
}

func callHttpbin(c context.Context) {
	n := rand.Intn(10)
	logrus.WithContext(c).Infof("Call httpbin with %d seconds delay", n)
	ctx := httptrace.WithClientTrace(c, otelhttptrace.NewClientTrace(c))
	resp, _ := otelhttp.Get(ctx, "http://httpbin.httpbin/delay/"+strconv.Itoa(n))
	defer resp.Body.Close()
}
