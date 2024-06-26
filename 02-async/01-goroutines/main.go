package main

import (
	"errors"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	go retry(func() error {
		return h.newsletterClient.AddToNewsletter(u)
	})

	go retry(func() error {
		return h.notificationsClient.SendNotification(u)
	})

	return nil
}

func retry(f func() error) {
	for {
		if err := f(); err != nil {
			continue
		}
		break
	}
}

func main() {
	var watermillLogger = log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))
	watermillLogger.Error("creating new redis stream publisher", errors.New("some err"), watermill.LogFields{})
}
