# CQRS with Consumer Groups

In both `NewEventHandler` and the custom `EventHandler`, you need to specify the handler name.
You can use the handler name to generate a consumer group name.
Thanks to that, you can have multiple handlers listening for the same event.

Remember the `SubscriberConstructor` option in the event processor config?
Typically, this option is used to dynamically generate a consumer group for each handler.
`SubscriberConstructor` is called for each handler.
Thanks to this, you no longer need to manually create a subscriber for each consumer group.

```go
return cqrs.NewEventProcessorWithConfig(
	cqrs.EventProcessorConfig{
		// ... 
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        redisClient, 
				ConsumerGroup: "svc-tickets." + params.HandlerName,
			}, watermillLogger)
		},
		// ...
    },
)
```

Consumer groups were explained in depth in [the previous exercise](/trainings/go-event-driven/exercise/f2eca145-e8cd-49c5-bdfa-384d33aa0bea).
If you need a refresher, check it out.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Watch out to not accidentally change the handler name.

Changing the handler name will lead to the creation of a new consumer group,
which may result in some old messages not being processed.

</span>
	</div>
	</div>

### Exercise

File: `09-cqrs-events/06-cqrs-with-consumer-groups/main.go`

Implement the `NewEventProcessor` function that will return a `*cqrs.EventProcessor` with a RedisStream consumer group generated based on the handler name.

The expected signature is:

```go
func NewEventProcessor(
	rdb *redis.Client, 
	router *message.Router,
    marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) (*cqrs.EventProcessor, error) {
```

Use the following for `GenerateSubscribeTopic`: 

```go
func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
	return params.EventName, nil
}
```

The consumer group name doesn't matter. The only requirement is that it should be different for each event handler. 

A common approach is to include the service name in it to avoid handler name conflicts across services.
For example, `"svc-users." + handlerName`.
