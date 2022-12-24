package services

import (
	"testing"
	"time"
)

const streamUrl = "http://live-test-scores.herokuapp.com/scores"

func TestExamService_SubscribeSSE(t *testing.T) {
	service := NewExamService()

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
			if err := service.SubscribeSSE(tt.sseUrl); (err != nil) != tt.wantErr {
				t.Errorf("ExamService.SubscribeSSE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExamService_ProcessMessage(t *testing.T) {
	service := NewExamService()

	tests := []struct {
		name    string
		message []byte
		wantErr bool
	}{
		// Failure
		{
			"nil message",
			nil,
			true,
		},
		{
			"empty message bytes",
			[]byte{},
			true,
		},
		{
			"invalid json content",
			[]byte("{123}"),
			true,
		},
		{
			"invalid message schema",
			[]byte(`{"studentNumber":"Marvin32","examId":11295,"examScore":0.6405682932226593}`),
			true,
		},
		{
			"empty json content",
			[]byte("{}"),
			true,
		},
		// Success
		{
			"valid message",
			[]byte(`{"studentId":"Marvin32","exam":11295,"score":0.6405682932226593}`),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.ProcessMessage(tt.message); (err != nil) != tt.wantErr {
				t.Errorf("ExamService.ProcessMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExamService_GetStudents(t *testing.T) {
	service := NewExamService()
	service.SubscribeSSE(streamUrl)
	defer service.StopSSESubscription()

	tests := []struct {
		name              string
		wantStudentsCount int
	}{
		{"no students", 0},
		{"valid students", 20},
	}

	for count, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStudents := service.GetStudents()

			// wait for the SSE stream to send data for the first time
			if count > 0 && len(gotStudents) == 0 {
				time.Sleep(5 * time.Second)
				gotStudents = service.GetStudents()
			}

			if len(gotStudents) != tt.wantStudentsCount {
				t.Errorf("ExamService.GetStudents() = %d, want %d", len(gotStudents), tt.wantStudentsCount)
			}
		})
	}
}

func TestExamService_GetExams(t *testing.T) {
	service := NewExamService()
	service.SubscribeSSE(streamUrl)
	defer service.StopSSESubscription()

	tests := []struct {
		name      string
		wantExams bool
	}{
		{"no exams", false},
		{"valid exams", true},
	}

	for count, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExams := service.GetExams()

			// wait for the SSE stream to send data for the first time
			if count > 0 && len(gotExams) == 0 {
				time.Sleep(5 * time.Second)
				gotExams = service.GetExams()
			}

			if (len(gotExams) == 0) == tt.wantExams {
				t.Errorf("ExamService.GetExams() = %d, want > 0", len(gotExams))
			}
		})
	}
}
