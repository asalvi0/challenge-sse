package tests

import (
	"testing"
	"time"

	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"
)

func TestExamController_GetExamNumbers(t *testing.T) {
	eventsCh := make(chan models.Event)
	eventController := controllers.NewEventController()
	eventController.StartSSESubscription(streamUrl, eventsCh)
	defer eventController.StopSSESubscription()

	examController := controllers.NewExamController(eventsCh)

	tests := []struct {
		name      string
		wantExams bool
	}{
		{"no exams", false},
		{"valid exams", true},
	}

	for count, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExams := examController.GetExamNumbers()

			// wait for the SSE stream to send data for the first time
			if count > 0 && len(gotExams) == 0 {
				time.Sleep(5 * time.Second)
				gotExams = examController.GetExamNumbers()
			}

			if (len(gotExams) == 0) == tt.wantExams {
				t.Errorf("ExamService.GetExams() = %d, want > 0", len(gotExams))
			}
		})
	}
}
