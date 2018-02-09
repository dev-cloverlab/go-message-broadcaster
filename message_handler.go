package broadcaster

import "context"

type MessageHandlerID int
type MessageHandler func(msg *RequestMessage, ctx context.Context) (*ResponseMessage, error)
type MessageHandlers map[MessageHandlerID]MessageHandler
