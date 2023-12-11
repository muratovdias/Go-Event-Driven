# Handling the Correlation ID

In the previous implementation, we were [manually setting](/trainings/go-event-driven/exercise/2d4b8c79-e619-4d11-bebe-e367d0d066e6) the `"correlation_id"` metadata in our message header.

There are at least two ways to handle this within the event bus:
- We can modify the marshaler.
- We can implement a publisher decorator.

### Publisher Decorator

The decorator pattern is not very popular in Go, which is unfortunate. 
It's a very powerful pattern that allows you to modify an object's behavior without changing its implementation.

An example of the decorator implementation is as follows:

```go
type Sender interface {
	Send(message string) error
}

type QuotingDecorator struct {
	Sender
}

func (q QuotingDecorator) Send(message string) error {
	message = fmt.Sprintf(`"%s"`, message)

	return q.Sender.Send(message)
}
```

As you can see, the decorator implements the same interface as the decorated object. 
It can modify arguments, return values, or even call the decorated object multiple times.

### Further Reading

If you want to learn more about the use cases of the decorator pattern, 
you should check out the [Increasing Cohesion in Go with Generic Decorators](https://threedots.tech/post/increasing-cohesion-in-go-with-generic-decorators/) article on our blog!

## Exercise

File: `09-cqrs-events/03-handling-correlation-id/main.go`

Let's implement a decorator for the publisher. 
This decorator should add the correlation ID to the message metadata key `correlation_id`. 
The decorator should be called `CorrelationPublisherDecorator`.

You can retrieve the correlation ID from the context using the `CorrelationIDFromContext` function.

Remember to set the metadata for each incoming message! Each message has its own context and metadata.
