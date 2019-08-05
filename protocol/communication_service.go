package protocol

// CommunicationService 数据传输
type CommunicationService interface {
	Send(msgID uint16, msgData []byte) (int, error)
	Parse() ([]byte, error)
}
