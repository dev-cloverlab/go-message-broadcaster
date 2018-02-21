package broadcaster

import "context"

type MessageHandlerID int
type MessageHandler func(msg *RequestMessage, ctx context.Context) (ResponseMessages, error)
type MessageHandlers map[MessageHandlerID]MessageHandler
