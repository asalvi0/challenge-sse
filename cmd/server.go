package main

import (
	"githb.com/asalvi0/challenge-sse/internal/api"
	"githb.com/asalvi0/challenge-sse/internal/services"
)

// TODO: read this value from a config/secrets store or env variable
const (
	streamUrl = "http://live-test-scores.herokuapp.com/scores"
	apiPort   = 8080
)

func main() {
	service := services.NewExamService()

	// subscribe to SSE stream
	service.SubscribeSSE(streamUrl)

	// start api server
	server := api.NewServer(apiPort, service)
	server.Start()
}
