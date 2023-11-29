package v1

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type publisherI interface {
	Publish(topic string, msgs ...*message.Message) error
	Close() error
}

type loggerI interface {
	Error(msg string, err error, fields watermill.LogFields)
	Info(msg string, fields watermill.LogFields)
	Debug(msg string, fields watermill.LogFields)
	Trace(msg string, fields watermill.LogFields)
	With(fields watermill.LogFields) watermill.LoggerAdapter
}
