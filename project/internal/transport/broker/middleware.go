package broker

import (
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		logger := log.FromContext(msg.Context())
		logger = logger.WithField("message_uuid", msg.UUID)
		logger.Info("Handling a message")

		defer func() {
			if err != nil {
				logger.WithFields(logrus.Fields{
					"error":        err,
					"message_uuid": msg.UUID,
				}).Error("Message handling error")
			}
		}()

		return next(msg)
	}
}

func PropagateCorrelationID(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")

		if correlationID == "" {
			correlationID = watermill.NewUUID()
		}

		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		ctx = log.ToContext(ctx, logrus.WithFields(logrus.Fields{"correlation_id": correlationID}))

		msg.SetContext(ctx)

		return next(msg)
	}
}
