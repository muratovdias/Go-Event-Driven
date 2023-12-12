package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type FollowRequestSent struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type EventsCounter interface {
	CountEvent() error
}

type EventHandler struct {
	counter EventsCounter
}

func (e *EventHandler) CountEvent(ctx context.Context, event *FollowRequestSent) error {
	return e.counter.CountEvent()
}

func NewFollowRequestSentHandler(counter EventsCounter) cqrs.EventHandler {

	h := EventHandler{
		counter: counter,
	}

	return cqrs.NewEventHandler(
		"CountEvent",
		h.CountEvent,
	)
}
