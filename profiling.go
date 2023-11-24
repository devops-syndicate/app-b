package main

import (
	"os"

	"github.com/grafana/pyroscope-go"
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
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
