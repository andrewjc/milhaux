package common

import "encoding/json"

type StorageConnector interface {
	PerformSendReceive(StorageMessage) (StorageMessage, error)
}

type StorageMessageType string

const (
	StorageMessageTypeQuery           = StorageMessageType("QUERY")
	StorageMessageTypeSuccessResponse = StorageMessageType("RESPONSE")
)

type StorageMessageCommand string

const (
	StorageMessageCommandNone    = StorageMessageCommand("NONE")
	StorageMessageCommandIsValid = StorageMessageCommand("ISVALID")
	StorageMessageCommandList    = StorageMessageCommand("LIST")
	StorageMessageCommandRename  = StorageMessageCommand("RENAME")
	StorageMessageCommandCreate  = StorageMessageCommand("CREATE")
)

type StorageMessage interface {
	MsgType() StorageMessageType
	MsgCommand() StorageMessageCommand
	MsgData() interface{}
	ToJson() string
}

type storageMessageImpl struct {
	msgType    StorageMessageType
	msgCommand StorageMessageCommand
	msgData    interface{}
}

func (m *storageMessageImpl) MsgType() StorageMessageType {
	return m.msgType
}

func (m *storageMessageImpl) MsgCommand() StorageMessageCommand {
	return m.msgCommand
}

func (m *storageMessageImpl) MsgData() interface{} {
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

func (s *StorageMessageBuilder) StorageCommandMessage(command string, args string) *StorageMessageBuilder {
	s.msgType = StorageMessageTypeQuery
	s.msgCommand = StorageMessageCommand(command)
	s.msgData = args

	return s
}

func (s *StorageMessageBuilder) SuccessResponseMessage(msgData interface{}) *StorageMessageBuilder {
	s.msgType = StorageMessageTypeSuccessResponse
	s.msgCommand = StorageMessageCommandNone
	s.msgData = msgData

	return s
}
