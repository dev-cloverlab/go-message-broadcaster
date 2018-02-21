package broadcaster

type Conn interface {
	SendMessage(*ResponseMessage) error
	Receive() (*RequestMessage, error)
	Close() error
}
