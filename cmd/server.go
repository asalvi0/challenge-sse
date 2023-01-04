package main

import (
	"log"

	"github.com/asalvi0/challenge-sse/internal/api"
	"github.com/asalvi0/challenge-sse/internal/controllers"
)

// TODO: read config from a config/secrets store or env variable
const (
	streamUrl = "http://live-test-scores.herokuapp.com/scores"
	port      = 8080
)

func main() {
	eventController := controllers.NewEventController(streamUrl)

	// subscribe to SSE stream
	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		log.Fatal(err)
	}

	examController := controllers.NewExamController(eventsCh)
	studentController := controllers.NewStudentController(eventsCh)

	// start API server
	apiServer := api.NewServer(port, eventController, examController, studentController)
	apiServer.Start()
}
