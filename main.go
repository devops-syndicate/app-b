package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	logrus.WithContext(c.Request.Context()).Info("Call hello endpoint")
	c.String(http.StatusOK, "Hello '%s'", APP_NAME)
}

func RandomHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	logrus.WithContext(c.Request.Context()).Infof("Call random endpoint and wait for %d seconds", n)
	time.Sleep(time.Duration(n) * time.Second)
	c.String(http.StatusOK, "From random")
}
