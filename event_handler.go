package broadcaster

import "context"

type EventHandler func(msg *EventMessage, ctx context.Context) (ResponseMessages, error)
type EventHandlers map[EventType]EventHandler
