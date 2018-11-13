package backend

import (
	"github.com/milhaux/common"
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
	packageConfig.backend = FILESTORE
	packageConfig.available = true
}

type MailStoreBackend interface {
	Start() error
}

func NewMailStoreBackend(config *common.ApplicationConfig) MailStoreBackend {
	log.Debug("Creating a mailstore backend instance...")

	switch packageConfig.backend {
	case MEMSTORE:
		return &MemStoreStorageBackend{config}
	case FILESTORE:
		return &FsStoreStorageBackend{config}
	case DBSTORE:
		return &DbStoreStorageBackend{config}
	}

	return nil
}
