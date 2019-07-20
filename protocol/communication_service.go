package protocol

// 数据传输
type CommunicationService interface {
	Send(data []byte) (int, error)
	Parse() ([]byte, error)
}
