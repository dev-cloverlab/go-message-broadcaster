package broadcaster

import "context"

type EventHandler func(msg *EventMessage, ctx context.Context) (*ResponseMessage, error)
type EventHandlers map[EventType]EventHandler
