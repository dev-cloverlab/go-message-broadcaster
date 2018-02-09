package broadcaster

type CastType int

const (
	All CastType = iota + 1
	Self
	Exclusive
)

type RequestMessage struct {
	SenderID  ClientID
	HandlerID MessageHandlerID
	Body      []byte
}

type ResponseMessage struct {
	SenderID  ClientID
	HandlerID MessageHandlerID
	CastType  CastType
	CastFor   []ClientID
	Body      []byte
}

func NewResponseMessage(ct CastType, body []byte, cf ...ClientID) *ResponseMessage {
	return &ResponseMessage{
		CastType: ct,
		CastFor:  cf,
		Body:     body,
	}
}

type EventType int

const (
	OnAddClient EventType = iota + 1
	OnDelClient
)

type EventMessage struct {
	EventType EventType
	ClientID  ClientID
}

func NewEventMessage(et EventType, cid ClientID) *EventMessage {
	return &EventMessage{
		EventType: et,
		ClientID:  cid,
	}
}
