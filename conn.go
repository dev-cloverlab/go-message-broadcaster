package broadcaster

type Conn interface {
	SendMessage(*ResponseMessage) error
	SendEvent(*EventMessage) error
	Receive() (*RequestMessage, error)
	Close() error
}
