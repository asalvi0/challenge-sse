package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/asalvi0/challenge-sse/internal/models"
	"github.com/r3labs/sse/v2"
)

type EventController struct {
	sseClient *sse.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewEventController() *EventController {
	controller := EventController{}
	controller.ctx, controller.ctxCancel = context.WithCancel(context.Background())

	return &controller
}

func validateSSEUrl(sseUrl string) (err error) {
	// validate url format
	parsedUrl, err := url.Parse(sseUrl)
	if (err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") || parsedUrl.Host == "") || len(sseUrl) == 0 {
		return errors.New("invalid SSE url")
	}

	// check if url is reachable
	// TODO: implement a proper validation method
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}

	resp, err := client.Get(sseUrl)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("unreachable SSE url")
	}

	return nil
}

func (c *EventController) StartSSESubscription(url string, eventsCh chan models.Event) (err error) {
	err = validateSSEUrl(url)
	if err != nil {
		return err
	}

	c.sseClient = sse.NewClient(url)
	if c.sseClient == nil {
		return errors.New("failed to create SSE client")
	}

	go func() {
		c.sseClient.SubscribeRawWithContext(c.ctx, func(msg *sse.Event) {
			event, err := models.ParseEvent(msg.Data)
			if err != nil {
				log.Println(err)
			}

			if len(event.StudentId) == 0 || event.Number <= 0 {
				log.Println(errors.New("invalid message data"))
			}

			eventsCh <- event
		})
	}()

	return nil
}

func (c *EventController) StopSSESubscription() {
	c.ctxCancel()
}
