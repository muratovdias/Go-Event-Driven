# Implementing the Event Handler

You know how to configure an Event Bus and an Event Processor. Now it's time to implement the Event Handler.

There are two ways to do it.

### Implementing the `EventHandler` interface

```go
// EventHandler receives events defined by NewEvent and handles them with its Handle method.
// If using DDD, CommandHandler may modify and persist the aggregate.
// It can also invoke a process manager or a saga or just build a read model.
//
// In contrast to CommandHandler, every Event can have multiple EventHandlers.
//
// One instance of EventHandler is used during handling messages.
// When multiple events are delivered at the same time, the Handle method can be executed multiple times at the same time.
// Because of this, the Handle method needs to be thread safe!
type EventHandler interface {
	// HandlerName is the name used in message.Router while creating the handler.
	//
	// It will be also passed to EventsSubscriberConstructor.
	// It may be useful, for example, to create a consumer group for each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result in **reconsuming all messages** !!!
	HandlerName() string

	NewEvent() interface{}

	Handle(ctx context.Context, event interface{}) error
}
```

Example implementation:

```go
type ArticlePublished struct {
	ArticleID string `json:"article_id"`
}

type ArticlePublishedHandler struct {}

func (h *ArticlePublishedHandler) HandlerName() string {
	return "ArticlePublishedHandler"
}

func (h *ArticlePublishedHandler) NewEvent() interface{} {
	return &ArticlePublished{}
}

func (h *ArticlePublishedHandler) Handle(ctx context.Context, event any) error {
	e := event.(*ArticlePublished)

	fmt.Printf("Article %s published\n", e.ArticleID)
	
	return nil
}
```

Based on the type returned from `NewEvent()`, the Event Processor will dispatch the received event to the proper handler.

### Using the generic `NewEventHandler` function

Recently, we've added support for the `NewEventHandler` function that uses generics under the hood.
It generates an `EventHandler` implementation dynamically based on a function signature: 

```go
cqrs.NewEventHandler(
	"ArticlePublishedHandler", 
	func(ctx context.Context, event *ArticlePublished) error {
		fmt.Printf("Article %s published\n", event.ArticleID)
		
		return nil
	},
),
```

It requires a lot less boilerplate to add handlers.

#### Injecting dependencies into handler from `NewEventHandler`

You may wonder if it's possible to inject dependencies to the handler using the generic approach.
You can use the same technique that you likely already use with HTTP handlers.
Create a struct that holds all dependencies, and then pass the method as a value to `cqrs.NewEventHandler`:

```go
type ArticlesHandler struct {
	notificationsService NotificationsService
}

func (h ArticlesHandler) PrintIDOnArticlePublished(ctx context.Context, event *ArticlePublished) error {
	fmt.Printf("Article %s published\n", event.ArticleID)
	
	return nil
}

func (h ArticlesHandler) NotifyUserOnArticlePublished(ctx context.Context, event *ArticlePublished) error {
	h.notificationsService.NotifyUser(event.ArticleID)
	
	return nil
}

func NewArticlesHandlers(notificationsService NotificationsService) []cqrs.EventHandler {
	h := ArticlesHandler{
		notificationsService: notificationsService,
	}

	return []cqrs.EventHandler{
		cqrs.NewEventHandler(
			"PrintIDOnArticlePublished", 
			h.PrintIDOnArticlePublished,
		), 
		cqrs.NewEventHandler(
			"NotifyUserOnArticlePublished", 
			h.NotifyUserOnArticlePublished, 
		),
	}
}
```

Note that the `ArticlesHandler` can have multiple handlers for the same event.

## Exercise

File: `09-cqrs-events/05-event-handler/main.go`

Implement the `NewFollowRequestSentHandler` function that returns `cqrs.EventHandler`.
It should accept `EventsCounter` as a parameter.

You can use `NewEventHandler` or implement the `EventHandler` interface.
`CountEvent()` should be called each time the event is handled.
