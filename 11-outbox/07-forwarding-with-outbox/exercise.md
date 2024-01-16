# Forwarding messages

Now it's time to practice the last piece of the puzzle: forwarding messages from Postgres to Redis Pub/Sub.

## Exercise

File: `11-outbox/07-forwarding-with-outbox/main.go`

Implement the `RunForwarder` function:

```go
func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {
```

The function should:

1. Create a new Postgres Subscriber (like in the [11-outbox/05-subscribing-from-sql](/trainings/go-event-driven/exercise/f06083e3-e1ba-4682-ae4d-3f007bab0c20) exercise).
2. Call `SubscribeInitialize` with `outboxTopic` (like in the [11-outbox/05-subscribing-from-sql](/trainings/go-event-driven/exercise/f06083e3-e1ba-4682-ae4d-3f007bab0c20) exercise).
3. Create a new Redis Publisher.
4. Create a new `forwarder.NewForwarder` with ForwarderTopic config set to `outboxTopic`.
5. Run the forwarder.

The function should be non-blocking, but it should wait for the forwarder to start.
To achieve this, you can use:

```go
go func() {
    err := fwd.Run(context.Background())
    if err != nil {
        panic(err)
    }
}()

<-fwd.Running()
```

