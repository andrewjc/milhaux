package main

import (
	"os"
	"runtime"

	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/core"
	log "github.com/sirupsen/logrus"

	"flag"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	appConfig := common.GetAppConfig()

	serverCore := core.NewApplicationContext(appConfig)

	if *flag.Bool("backend_only", false, "a bool") {
		serverCore.Config.GetSmtpServerConfig().Enabled = false
		serverCore.Config.GetImap4ServerConfig().Enabled = false

		serverCore.Config.GetBackendConfig().ListeningInterface
	}

	serverCore.Start()

}
