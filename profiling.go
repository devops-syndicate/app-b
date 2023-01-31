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
		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
