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
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	logrus.WithContext(c).Infof("Call httpbin with %d seconds delay", n)
	ctx := httptrace.WithClientTrace(c, otelhttptrace.NewClientTrace(c))
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/delay/"+strconv.Itoa(n), nil)
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client.Do(req)
}
