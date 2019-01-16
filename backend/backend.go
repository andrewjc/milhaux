package backend

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/smtp"
	log "github.com/sirupsen/logrus"
)

var packageConfig = DefaultBackendConfig{}

type DefaultBackendConfig struct {
	backend   int8
	available bool
}

const (
	MEMSTORE = iota
	FILESTORE
	DBSTORE
)

// The package's init method
func init() {
	log.Debug("Init mailstore package")

	// Override config values for testing.
	packageConfig.backend = MEMSTORE
	packageConfig.available = true
}

type MailStoreBackend struct {
	smtpComponentChannel chan common.MailMessage
}

type MailStoreBackendProvider interface {
	Start() error
	IsStarted() bool
	OnSubmitQueue(message *smtp.SmtpServerChannelMessage)
}

func NewMailStoreBackend(config *common.ApplicationConfig) MailStoreBackendProvider {
	log.Debug("Creating a mailstore backend instance...")

	switch {
	case packageConfig.backend == MEMSTORE:
		return &MemStoreStorageBackend{}
	case packageConfig.backend == FILESTORE:
		return &FsStoreStorageBackend{}
	case packageConfig.backend == DBSTORE:
		return &DbStoreStorageBackend{}
	}

	return nil
}
