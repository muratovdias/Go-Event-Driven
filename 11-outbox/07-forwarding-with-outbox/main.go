package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {
	// your code goes here
	sub, err := watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
			OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	err = sub.SubscribeInitialize(outboxTopic)
	if err != nil {
		return err
	}

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		logger.Error("creating new redis stream publisher", err, watermill.LogFields{})
		panic(err)
	}

	forward, err := forwarder.NewForwarder(
		sub, publisher, logger,
		forwarder.Config{
			ForwarderTopic: outboxTopic,
		},
	)
	if err != nil {
		logger.Error("creating new forwarder", err, watermill.LogFields{})
		panic(err)
	}

	go func() {
		if err := forward.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

	forward.Running()

	return nil
}
