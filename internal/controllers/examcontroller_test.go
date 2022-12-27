package controllers

import (
	"testing"
	"time"

	"github.com/asalvi0/challenge-sse/internal/models"
)

func TestExamController_AddExam(t *testing.T) {
	eventsCh := make(chan models.Event)
	eventController := NewEventController()
	eventController.StartSSESubscription(streamUrl, eventsCh)
	defer eventController.StopSSESubscription()

	examController := NewExamController(eventsCh)

	tests := []struct {
		name    string
		exams   models.ExamsRepository
		event   models.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ExamController{
				exams:    tt.exams,
				eventsCh: tt.eventsCh,
			}
			if err := c.AddExam(tt.event); (err != nil) != tt.wantErr {
				t.Errorf("ExamController.AddExam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExamController_GetExamNumbers(t *testing.T) {
	eventsCh := make(chan models.Event)
	eventController := NewEventController()
	eventController.StartSSESubscription(streamUrl, eventsCh)
	defer eventController.StopSSESubscription()

	examController := NewExamController(eventsCh)

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
