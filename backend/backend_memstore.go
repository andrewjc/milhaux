package backend

import (
	linkedlist "container/list"
	"github.com/andrewjc/milhaux/common"
	log "github.com/sirupsen/logrus"
	"sync"
)

// The package's init method
func init() {
	log.Debug("Init mem storage backend package")
}

type MemStoreStorageBackend struct {
	config    *common.ApplicationConfig
	queue     *linkedlist.List
	isReady   bool
	initMutex sync.Mutex
}

var _isStartedTVal bool

func (s *MemStoreStorageBackend) IsStarted() bool {

	if _isStartedTVal {
		return true
	} //avoid unnecessary locking once initialized

	s.initMutex.Lock()
	_isStartedTVal = s.isReady
	s.initMutex.Unlock()
	return _isStartedTVal
}

func (s *MemStoreStorageBackend) Start() error {
	log.Info("Starting memory backed mailstore backend instance...")

	s.initMutex.Lock()
	s.queue = linkedlist.New()
	s.isReady = true
	s.initMutex.Unlock()

	return nil
}

func (s *MemStoreStorageBackend) Store(message *common.MailMessage) error {
	log.Infof("Added message to memstore. There are now %v items in the store.", s.queue.Len())
	s.queue.PushBack(message)
	return nil
}
