package broker

import (
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func (b *broker) initEventProcessor() {
	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		b.router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        b.rdb,
					ConsumerGroup: "svc-tockets" + params.HandlerName,
				}, b.watermillLogger)
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
		})
	if err != nil {
		panic(err)
	}

	// assign event processor
	b.eventProcessor = eventProcessor
}
