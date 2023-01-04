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
	streamUrl string
	sseClient *sse.Client
	eventsCh  chan models.Event
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewEventController(url string) *EventController {
	controller := EventController{streamUrl: url}

	controller.ctx, controller.ctxCancel = context.WithCancel(context.Background())
	controller.eventsCh = make(chan models.Event)

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

func (c *EventController) StopSSESubscription() {
	c.ctxCancel()
}

func (c *EventController) StartSSESubscription() (eventsCh chan models.Event, err error) {
	err = validateSSEUrl(c.streamUrl)
	if err != nil {
		return nil, err
	}

	c.sseClient = sse.NewClient(c.streamUrl)
	if c.sseClient == nil {
		return nil, errors.New("failed to create SSE client")
	}

	go func() {
		c.sseClient.SubscribeRawWithContext(c.ctx, func(msg *sse.Event) {
			event, err := models.ParseEvent(msg.Data)
			if err != nil {
				log.Println(err)
				return
			}

			if len(event.StudentId) == 0 || event.Number <= 0 {
				log.Println(errors.New("invalid message data"))
				return
			}

			eventsCh <- event
		})
	}()

	return c.eventsCh, nil
}
