package main

import (
	"encoding/json"
	"log"

	"githb.com/asalvi0/challenge-sse/internal/api"
	"githb.com/asalvi0/challenge-sse/internal/model"

	"github.com/r3labs/sse/v2"
)

// TODO: read this value from a config/secrets store or env variable
const streamUrl = "http://live-test-scores.herokuapp.com/scores"

func main() {
	// start api server
	server := api.NewServer(8080)
	go server.Start()

	// subscribe to SSE stream
	sseClient := sse.NewClient(streamUrl)
	sseClient.SubscribeRaw(func(msg *sse.Event) {
		var exam model.Data

		err := json.Unmarshal([]byte(msg.Data), &exam)
		if err != nil {
			log.Println(err)
			return
		}

		// TODO: process data
		log.Println(exam)
	})
}
