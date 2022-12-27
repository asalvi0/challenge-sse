package tests

import (
	"net/http"
	"testing"

	"github.com/asalvi0/challenge-sse/internal/api"
	"github.com/julienschmidt/httprouter"
)

func TestServer_Start(t *testing.T) {
	tests := []struct {
		name string
		s    *api.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Start()
		})
	}
}

func TestServer_getStudents(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		s    *api.Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.getStudents(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestServer_getStudent(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		params httprouter.Params
	}
	tests := []struct {
		name string
		s    *api.Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.getStudent(tt.args.w, tt.args.r, tt.args.params)
		})
	}
}

func TestServer_getExams(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		s    *api.Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.getExams(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestServer_getExam(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		params httprouter.Params
	}
	tests := []struct {
		name string
		s    *api.Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.getExam(tt.args.w, tt.args.r, tt.args.params)
		})
	}
}
