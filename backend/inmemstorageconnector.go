package backend

import "github.com/andrewjc/milhaux/common"

/*
	A simple backend storage connector that doesn't perform
	any serialization/deserialization and just passes the message
	directly to the backend.
*/

type InMemoryStorageConnector struct {
	MailStoreBackend
}

func (m *InMemoryStorageConnector) PerformSendReceive(message common.StorageMessage) (common.StorageMessage, error) {
	msg, err := m.handleConnectorMessage(message)
	return msg, err
}
