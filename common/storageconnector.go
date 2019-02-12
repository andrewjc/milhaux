package common

import "encoding/json"

type StorageConnector interface {
	PerformSendReceive(StorageMessage) (StorageMessage, error)
}

type StorageMessageType string

const (
	StorageMessageTypeQuery    = StorageMessageType("QUERY")
	StorageMessageTypeResponse = StorageMessageType("RESPONSE")
)

type StorageMessageCommand string

const (
	StorageMessageCommandNone    = StorageMessageCommand("NONE")
	StorageMessageCommandIsValid = StorageMessageCommand("ISVALID")
)

type StorageMessage interface {
	MsgType() StorageMessageType
	MsgCommand() StorageMessageCommand
	MsgData() string
	ToJson() string
}

type storageMessageImpl struct {
	msgType    StorageMessageType
	msgCommand StorageMessageCommand
	msgData    string
}

func (m *storageMessageImpl) MsgType() StorageMessageType {
	return m.msgType
}

func (m *storageMessageImpl) MsgCommand() StorageMessageCommand {
	return m.msgCommand
}

func (m *storageMessageImpl) MsgData() string {
	return m.msgData
}

func (m *storageMessageImpl) ToJson() string {
	ret, _ := json.Marshal(m)
	return string(ret)
}

type StorageMessageBuilder struct {
	storageMessageImpl
}

func (s *StorageMessageBuilder) Build() StorageMessage {
	return &s.storageMessageImpl
}

func (s *StorageMessageBuilder) IsValidCommandMessage(command string) *StorageMessageBuilder {
	s.msgType = StorageMessageTypeQuery
	s.msgCommand = StorageMessageCommandIsValid
	s.msgData = command

	return s
}

func (s *StorageMessageBuilder) ResponseMessage(msgData string) *StorageMessageBuilder {
	s.msgType = StorageMessageTypeResponse
	s.msgCommand = StorageMessageCommandNone
	s.msgData = msgData

	return s
}
