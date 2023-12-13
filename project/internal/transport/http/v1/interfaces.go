package v1

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
)

type publisherI interface {
	Publish(ctx context.Context, event any) error
}

type loggerI interface {
	Error(msg string, err error, fields watermill.LogFields)
	Info(msg string, fields watermill.LogFields)
	Debug(msg string, fields watermill.LogFields)
	Trace(msg string, fields watermill.LogFields)
	With(fields watermill.LogFields) watermill.LoggerAdapter
}
