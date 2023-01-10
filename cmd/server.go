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

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		log.Fatal(err)
	}

	examController, err := controllers.NewExamController(eventsCh)
	if err != nil {
		log.Fatal(err)
	}

	studentController, err := controllers.NewStudentController(eventsCh)
	if err != nil {
		log.Fatal(err)
	}

	apiServer := api.NewServer(port, eventController, examController, studentController)
	apiServer.Start()
}
