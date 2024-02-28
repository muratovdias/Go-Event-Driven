package command

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewCommandBus(pub message.Publisher) (*cqrs.CommandBus, error) {
	config := cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil
		},
		Marshaler: marshaller,
	}

	return cqrs.NewCommandBusWithConfig(pub, config)
}
