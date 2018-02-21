package broadcaster

type Conn interface {
	Send(*ResponseMessage) error
	Receive() (*RequestMessage, error)
	Close() error
}
