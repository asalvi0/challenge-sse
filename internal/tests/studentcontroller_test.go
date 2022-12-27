package tests

import (
	"testing"
	"time"

	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"
)

func TestStudentController_GetStudentsID(t *testing.T) {
	eventsCh := make(chan models.Event)
	eventController := controllers.NewEventController()
	eventController.StartSSESubscription(streamUrl, eventsCh)
	defer eventController.StopSSESubscription()

	studentController := controllers.NewStudentController(eventsCh)

	tests := []struct {
		name              string
		wantStudentsCount int
	}{
		{"no students", 0},
		{"valid students", 20},
	}

	for count, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStudents := studentController.GetStudentsID()

			// wait for the SSE stream to send data for the first time
			if count > 0 && len(gotStudents) == 0 {
				time.Sleep(5 * time.Second)
				gotStudents = studentController.GetStudentsID()
			}

			if len(gotStudents) != tt.wantStudentsCount {
				t.Errorf("ExamService.GetStudents() = %d, want %d", len(gotStudents), tt.wantStudentsCount)
			}
		})
	}
}
