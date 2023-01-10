package controllers

import (
	"testing"
)

const streamUrl = "http://live-test-scores.herokuapp.com/scores"

func TestEventController_StartSSESubscription(t *testing.T) {
	tests := []struct {
		name       string
		controller *EventController
		wantErr    bool
	}{
		// Failure
		{"invalid SSE url", NewEventController("http://"), true},
		{"empty SSE url", NewEventController(""), true},
		{"inexistent SSE url", NewEventController("http://asd.asd/sse"), true},
		// Success
		{"valid SSE url", NewEventController("http://live-test-scores.herokuapp.com/scores"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.controller.StartSSESubscription()
			if (err != nil) != tt.wantErr {
				t.Errorf("EventController.StartSSESubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
