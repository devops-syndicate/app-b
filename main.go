package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const APP_NAME = "app-b"

func main() {
	initLogger()
	initTracing()
	initProfiling()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(otelgin.Middleware(APP_NAME))

	router.GET("/", HelloHandler)

	logrus.Info("Start listening on port 8080")

	// Start and run the server
	router.Run(":8080")
}

func HelloHandler(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("Call hello endpoint")
	c.String(http.StatusOK, "Hello '%s'", APP_NAME)
}
