package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Header struct {
	ID          string `json:"id"`
	EventName   string `json:"event_name"`
	PublishedAt string `json:"occurred_at"`
}

type ProductOutOfStock struct {
	Header    Header `json:"header"`
	ProductID string `json:"product_id"`
}

type ProductBackInStock struct {
	Header    Header `json:"header"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) *Publisher {
	return &Publisher{
		pub: pub,
	}
}

func NewHeader(eventName string) Header {
	return Header{
		ID:          uuid.NewString(),
		EventName:   eventName,
		PublishedAt: time.Now().Format(time.RFC3339),
	}
}

func (p *Publisher) PublishProductOutOfStock(productID string) error {
	header := NewHeader("ProductOutOfStock")

	event := ProductOutOfStock{
		Header:    header,
		ProductID: productID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}

func (p *Publisher) PublishProductBackInStock(productID string, quantity int) error {
	header := NewHeader("ProductBackInStock")

	event := ProductBackInStock{
		Header:    header,
		ProductID: productID,
		Quantity:  quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}
