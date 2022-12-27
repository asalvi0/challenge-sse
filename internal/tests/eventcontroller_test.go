package tests

import (
	"testing"

	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"
)

const streamUrl = "http://live-test-scores.herokuapp.com/scores"

func TestExamService_StartSSESubscription(t *testing.T) {
	eventsCh := make(chan models.Event)
	eventController := controllers.NewEventController()

	tests := []struct {
		name    string
		sseUrl  string
		wantErr bool
	}{
		// Failure
		{"invalid SSE url", "http://", true},
		{"empty SSE url", "", true},
		{"inexistent SSE url", "http://asd.asd/sse", true},
		// Success
		{"valid SSE url", "http://live-test-scores.herokuapp.com/scores", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := eventController.StartSSESubscription(tt.sseUrl, eventsCh); (err != nil) != tt.wantErr {
				t.Errorf("ExamService.StartSSESubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
