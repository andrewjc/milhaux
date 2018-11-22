package common

type MainEventMessage struct {
	data        interface{}
	MessageType interface{}
}

const (
	SHUTDOWN = iota
	PING
)

var (
	eventLoopChannel chan MainEventMessage
)

func GetMainMessageLoop() chan MainEventMessage {
	return eventLoopChannel
}
