package main

import (
	"github.com/asalvi0/challenge-sse/internal/api"
	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"
)

// TODO: read config from a config/secrets store or env variable
const (
	streamUrl = "http://live-test-scores.herokuapp.com/scores"
	port      = 8080
)

func main() {
	eventsCh := make(chan models.Event)

	eventController := controllers.NewEventController()
	examController := controllers.NewExamController(eventsCh)
	studentController := controllers.NewStudentController(eventsCh)

	// subscribe to SSE stream
	eventController.StartSSESubscription(streamUrl, eventsCh)

	// start API server
	apiServer := api.NewServer(port, eventController, examController, studentController)
	apiServer.Start()
}
