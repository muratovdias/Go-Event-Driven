package main

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) (*cqrs.EventBus, error) {

	config := cqrs.EventBusConfig{
		GeneratePublishTopic: GeneratePublishTopic,
		Marshaler:            &cqrs.JSONMarshaler{},
	}

	return cqrs.NewEventBusWithConfig(pub, config)
}

func GeneratePublishTopic(params cqrs.GenerateEventPublishTopicParams) (string, error) {
	return params.EventName, nil
}
