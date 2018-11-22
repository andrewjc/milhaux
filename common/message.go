package common

type MailMessageType uint8

type MailMessage struct {
	To   string
	From string
	Data string

	QueueId string
}
