package protocol

// 数据传输
type CommunicationService interface {
	Send(msgId uint16, msgData []byte) (int, error)
	Parse() ([]byte, error)
}
