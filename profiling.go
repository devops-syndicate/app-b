package main

import (
	"os"

	"github.com/pyroscope-io/client/pyroscope"
	"github.com/sirupsen/logrus"
)

func initProfiling() {

	var pyroscope_endpoint = os.Getenv("PYROSCOPE_URL")

	if pyroscope_endpoint == "" {
		logrus.Info("PYROSCOPE_URL not set. Skip profiling setup")
		return
	}

	pyroscope.Start(pyroscope.Config{
		ApplicationName: APP_NAME,
		ServerAddress:   pyroscope_endpoint,
		Logger:          logrus.StandardLogger(),
	})
}
