package main

import (
	"os"
	"runtime"

	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/core"
	log "github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	appConfig := common.GetAppConfig()

	serverCore := core.NewApplicationContext(appConfig)
	serverCore.Start()
}
