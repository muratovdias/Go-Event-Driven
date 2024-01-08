# Idempotency key

Another deduplication strategy is using an idempotency key (sometimes called a deduplication ID). 
This is a unique identifier of a request that is sent by the client. 
If the server receives a request with the same idempotency key, it will not process it again. 

One of the most common scenarios for using an idempotency key are payment service provider (PSP) APIs.
They can ensure that we are not charging the customer twice.

Let's see how this can work in practice.

## Exercise

File: `10-at-least-once-delivery/10-idempotency-key-in-api/main.go`

Modify the `PalPayClient` so  that _idempotency key_ will be sent in the `Idempotency-Key` header.
