package controllers

import (
	"testing"
	"time"

	"github.com/asalvi0/challenge-sse/internal/models"
)

func TestNewStudentController(t *testing.T) {
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
			controller, err := NewStudentController(tt.eventsCh)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewStudentController() error = %v, wantErr = %t", err, tt.wantErr)
				return
			}

			if (controller != nil) != tt.wantController {
				t.Errorf("NewStudentController() = %v, wantController = %t", controller, tt.wantController)
			}
		})
	}
}

func TestStudentController_AddStudent(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	studentController, err := NewStudentController(eventsCh)
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
			err := studentController.AddStudent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentController.AddStudent() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestStudentController_GetStudentsID(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	studentController, err := NewStudentController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		wantStudentsID bool
	}{
		{"no students", false},
		{"valid students", true},
	}

	for testIndex, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStudents := studentController.GetStudentsID()

			// wait for the SSE stream to send data
			if testIndex > 0 && len(gotStudents) == 0 {
				for {
					gotStudents = studentController.GetStudentsID()

					if len(gotStudents) != 0 {
						break
					}
					time.Sleep(1 * time.Second)
				}
			}

			if (len(gotStudents) > 0) != tt.wantStudentsID {
				t.Errorf("StudentController.GetStudentsID() = %v, wantStudentsID = %t", gotStudents, tt.wantStudentsID)
			}
		})
	}
}

func TestStudentController_GetStudent(t *testing.T) {
	eventController := NewEventController(streamUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

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
		id          string
		wantScores  bool
		wantAverage bool
	}{
		// Failure
		{"empty studentId", "", false, false},
		{"invalid studentId", "-1", false, false},

		// Success
		{"valid studentId", gotStudents[0].Id, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScores, gotAverage := studentController.GetStudent(tt.id)

			if (len(gotScores) > 0) != tt.wantScores {
				t.Errorf("StudentController.GetStudent() gotScores = %v, wantScores = %t", gotScores, tt.wantScores)
			}

			if (gotAverage > 0) != tt.wantAverage {
				t.Errorf("StudentController.GetStudent() gotAverage = %v, wantAverage = %t", gotAverage, tt.wantAverage)
			}
		})
	}
}
