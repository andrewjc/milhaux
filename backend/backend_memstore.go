package backend

import (
	linkedlist "container/list"
	"github.com/andrewjc/milhaux/common"
	. "github.com/andrewjc/milhaux/smtp"
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
func (backend *MemStoreStorageBackend) IsStarted() bool {

	if _isStartedTVal { return true } //avoid unnecessary locking once initialized

	backend.initMutex.Lock()
	_isStartedTVal = backend.isReady
	backend.initMutex.Unlock()
	return _isStartedTVal
}

const MAX_QUEUE_WORKERS = 4

func (backend *MemStoreStorageBackend) Start() error {
	log.Info("Starting memory backed mailstore backend instance...")

	backend.initMutex.Lock()
	backend.queue = linkedlist.New()
	backend.isReady = true
	backend.initMutex.Unlock()

	for i := 0; i < MAX_QUEUE_WORKERS; i++ {
		go backend.QueueWorker()
	}

	return nil
}

func (backend *MemStoreStorageBackend) OnSubmitQueue(message *SmtpServerChannelMessage) {
	backend.queue.PushBack(message)
}

func (backend *MemStoreStorageBackend) QueueWorker() {
	for {

	}
}
