# Writing Mocks

You need to be able to do all three of the following to execute component tests:
- Run the service as a function
- Mock external dependencies
- Inject dependencies that run in Docker

We will cover all of these in subsequent modules.

Let's start with writing mocks for our external dependencies. 
Remember that we are not counting Redis or Postgres as external dependencies — we can run them in Docker.

We should start with abstracting our external dependencies with interfaces.

```go
type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}
```

This is an example implementation of the interface above using clients from the common library:

```go
package api

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
)

type ReceiptsServiceClient struct {
	// We are not mocking this client: it's pointless to use the interface here
	clients *clients.Clients
}

func NewReceiptsServiceClient(clients *clients.Clients) *ReceiptsServiceClient {
	if clients == nil {
		panic("NewReceiptsServiceClient: clients is nil")
	}

	return &ReceiptsServiceClient{clients: clients}
}

func (c ReceiptsServiceClient) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, receipts.CreateReceipt{
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
		TicketId: request.TicketID,
	})
	if err != nil {
		return entities.IssueReceiptResponse{}, fmt.Errorf("failed to post receipt: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		// Receipt already exists
		return entities.IssueReceiptResponse{
			ReceiptNumber: resp.JSON200.Number,
			IssuedAt:      resp.JSON200.IssuedAt,
		}, nil
	case http.StatusCreated:
		// Receipt was created
		return entities.IssueReceiptResponse{
			ReceiptNumber: resp.JSON201.Number,
			IssuedAt:      resp.JSON201.IssuedAt,
		}, nil
	default:
		return entities.IssueReceiptResponse{}, fmt.Errorf("unexpected status code for POST receipts-api/receipts: %d", resp.StatusCode())
	}
}
```

Thanks to this, we can now create a mock implementation of `ReceiptsService`.

Note that we are not returning or accepting any types provided by the common library API clients except `entities.IssueReceiptRequest` and `entities.IssueReceiptResponse`.
It's a good practice to not propagate any external dependencies to the rest of the application; any breaking changes to our API client will force us to make changes in multiple places.
It's also a good idea to not expose all the data that we get from external APIs,  just the data that we need. 
This will decrease coupling and make our code simpler.


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

**We do not recommend using any libraries for mocking.**
They are usually very fragile and hard to debug.

Writing mocks by hand is usually quick and gives you a lot of flexibility.

</span>
	</div>
	</div>

This is how the mock implementation might look:

```go
type DataThatWeNeedToPass struct {
	TicketID string
}

type DataThatWeNeedToReturn struct {
	Number string
	DoneAt time.Time
}

type SomeMock struct {
	// let's make it thread-safe
	mock sync.Mutex

	PassedData []DataThatWeNeedToPass
}

func (c *SomeMock) DoStuff(ctx context.Context, request DataThatWeNeedToPass) (DataThatWeNeedToReturn, error) {
	c.mock.Lock()
	defer c.mock.Unlock()

	c.PassedData = append(c.PassedData, request)

	return DataThatWeNeedToReturn{
        Number: "mocked-number",
		DoneAt: time.Now(),
	}, nil
}
```

It's possible to write such a mock in a minute.

We have a lot of flexibility about how the mock should work.
We can also model any external logic.

The benefit compared to a generated mock is that we are writing mock logic only once, not for each test.
We also don't need to define which function calls we expect.

In tests, we can assert what data was provided by accessing the `PassedData` field.

This approach assumes that external services will behave in the way that you implemented in the mock.
It's not possible to ensure that this is always true, but it's also not a promise of any kind of tests.
Tests should give you more confidence that your code works as expected, but they are not guarantees.

## Exercise

File: `08-component-testing/03-writing-mocks/main.go`

Write a mock for the provided `ReceiptsService` interface.

The mock should be a struct named `ReceiptsServiceMock`.

It should have a public field `IssuedReceipts []IssueReceiptRequest` that stores all the requests that were passed to the `IssueReceipt` method.
`ReceiptsServiceMock.IssueReceipt` should return `IssueReceiptResponse` with non-zero fields.

`ReceiptsServiceMock` should be thread-safe — your mock will be used later in component tests that are executed in parallel.

To ensure that `ReceiptsServiceMock` is thread-safe, add locks like this:


```go
func (c *SomeMock) SomeMethod() {
	c.mock.Lock()
	defer c.mock.Unlock()
	
	// ...
```


<div class="alert alert-dismissible bg-light-primary d-flex flex-column flex-sm-row p-7 mb-10">
    <div class="d-flex flex-column">
        <h3 class="mb-5 text-dark">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-lightbulb text-primary" viewBox="0 0 16 16">
			  <path d="M2 6a6 6 0 1 1 10.174 4.31c-.203.196-.359.4-.453.619l-.762 1.769A.5.5 0 0 1 10.5 13a.5.5 0 0 1 0 1 .5.5 0 0 1 0 1l-.224.447a1 1 0 0 1-.894.553H6.618a1 1 0 0 1-.894-.553L5.5 15a.5.5 0 0 1 0-1 .5.5 0 0 1 0-1 .5.5 0 0 1-.46-.302l-.761-1.77a1.964 1.964 0 0 0-.453-.618A5.984 5.984 0 0 1 2 6zm6-5a5 5 0 0 0-3.479 8.592c.263.254.514.564.676.941L5.83 12h4.342l.632-1.467c.162-.377.413-.687.676-.941A5 5 0 0 0 8 1z"/>
			</svg>
			Tip
		</h3>
        <span>

Don't forget to use a pointer receiver for the `ReceiptsServiceMock.IssueReceipt` method.
Otherwise, it won't be able to modify any struct fields.
If your code doesn't work and there's no obvious reason, check that this is not the issue.

```go
func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error) {
```
