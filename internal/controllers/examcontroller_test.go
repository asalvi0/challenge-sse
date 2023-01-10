package controllers

import (
	"testing"
	"time"

	"github.com/asalvi0/challenge-sse/internal/models"
)

func TestNewExamController(t *testing.T) {
	tests := []struct {
		name           string
		eventsCh       <-chan models.Event
		wantController bool
		wantErr        bool
	}{
		// Failure
		{"nil channel", nil, false, true},

		// Success
		{"valid channel", make(<-chan models.Event), true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := NewExamController(tt.eventsCh)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewExamController() error = %v, wantErr = %t", err, tt.wantErr)
				return
			}

			if (controller != nil) != tt.wantController {
				t.Errorf("NewExamController() = %v, wantController = %t", controller, tt.wantController)
			}
		})
	}
}

func TestExamController_AddExam(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	examController, err := NewExamController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		event   models.Event
		wantErr bool
	}{
		// Failure
		{"empty exam result", models.Event{}, true},
		{
			"only studentId",
			models.Event{
				StudentId: "student-1",
			},
			true,
		},
		{
			"only empty studentId",
			models.Event{
				StudentId: "",
			},
			true,
		},
		{
			"no studentId",
			models.Event{
				ExamResult: models.ExamResult{
					Number: 1, Score: 0.7,
				},
			},
			true,
		},
		{
			"only score",
			models.Event{
				ExamResult: models.ExamResult{
					Score: 0.7,
				},
			},
			true,
		},
		{
			"only invalid score",
			models.Event{
				ExamResult: models.ExamResult{
					Score: -1,
				},
			},
			true,
		},
		{
			"only exam number",
			models.Event{
				ExamResult: models.ExamResult{
					Number: 1,
				},
			},
			true,
		},
		{
			"only invalid exam number",
			models.Event{
				ExamResult: models.ExamResult{
					Number: -1,
				},
			},
			true},
		{
			"no score",
			models.Event{
				StudentId: "student-1",
				ExamResult: models.ExamResult{
					Number: 1,
				},
			},
			true,
		},
		{
			"no exam number",
			models.Event{
				StudentId: "student-1",
				ExamResult: models.ExamResult{
					Score: 0.7,
				},
			},
			true,
		},

		// Success
		{
			"valid exam result",
			models.Event{
				StudentId: "student-1",
				ExamResult: models.ExamResult{
					Number: 1,
					Score:  0.7,
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := examController.AddExam(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExamController.AddExam() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestExamController_GetExamNumbers(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	examController, err := NewExamController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		wantExams bool
	}{
		{"no exams", false},
		{"valid exams", true},
	}

	for testIndex, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExams := examController.GetExamNumbers()

			// wait for the SSE stream to send data
			if testIndex > 0 && len(gotExams) == 0 {
				for {
					gotExams = examController.GetExamNumbers()

					if len(gotExams) != 0 {
						break
					}
					time.Sleep(1 * time.Second)
				}
			}

			if (len(gotExams) > 0) != tt.wantExams {
				t.Errorf("ExamController.GetExams() = %v, wantExams = %t", gotExams, tt.wantExams)
			}
		})
	}
}

func TestExamController_GetExam(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	examController, err := NewExamController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	studentController, err := NewStudentController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	gotStudents := studentController.GetStudents()

	// wait for the SSE stream to send data
	if len(gotStudents) == 0 {
		for {
			gotStudents = studentController.GetStudents()

			if len(gotStudents) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	tests := []struct {
		name        string
		number      int
		wantScores  bool
		wantAverage bool
	}{
		// Failure
		{"invalid exam number", -1, false, false},

		// Success
		{"valid exam number", gotStudents[0].Exams[0].Number, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScores, gotAverage := examController.GetExam(tt.number, gotStudents)

			if (len(gotScores) > 0) != tt.wantScores {
				t.Errorf("ExamController.GetExam() gotScores = %v, wantScores = %t", gotScores, tt.wantScores)
			}

			if (gotAverage > 0) != tt.wantAverage {
				t.Errorf("ExamController.GetExam() gotAverage = %v, wantAverage = %t", gotAverage, tt.wantAverage)
			}
		})
	}
}
