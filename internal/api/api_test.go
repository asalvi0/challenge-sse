package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"
	"github.com/julienschmidt/httprouter"
)

var sseUrl = "http://live-test-scores.herokuapp.com/scores"

func Router(server *Server) *httprouter.Router {
	router := httprouter.New()

	router.GET("/students", server.getStudents)
	router.GET("/students/:id", server.getStudent)

	router.GET("/exams", server.getExams)
	router.GET("/exams/:number", server.getExam)

	router.GET("/start-sse", server.startSSESubscription)
	router.GET("/stop-sse", server.stopSSESubscription)

	return router
}

func TestServer_getStudents(t *testing.T) {
	eventController := controllers.NewEventController(sseUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	studentController, err := controllers.NewStudentController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	// wait for the SSE stream to send data
	if len(studentController.GetStudentsID()) == 0 {
		for {
			if len(studentController.GetStudentsID()) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	tests := []struct {
		name           string
		server         *Server
		wantStatusCode int
		wantData       bool
	}{
		// Failure
		{
			"no student controller",
			NewServer(0, eventController, nil, nil),
			500,
			false,
		},

		// Success
		{
			"no event controller",
			NewServer(0, nil, nil, studentController),
			200,
			true,
		},
		{
			"no exam controller",
			NewServer(0, eventController, nil, studentController),
			200,
			true,
		},
		{
			"valid student controller",
			NewServer(0, eventController, nil, studentController),
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/students", nil)
			response := httptest.NewRecorder()

			Router(tt.server).ServeHTTP(response, request)
			resp := response.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Server.getStudents() gotStatusCode = %v, wantStatusCode = %v", resp.StatusCode, tt.wantStatusCode)
			}

			var gotData []interface{}
			_ = json.Unmarshal(body, &gotData) // intentionally ignore error

			if (len(gotData) != 0) != tt.wantData {
				t.Errorf("Server.getStudents() gotData = %v, wantData = %v", len(gotData) != 0, tt.wantData)
			}
		})
	}
}

func TestServer_getStudent(t *testing.T) {
	eventController := controllers.NewEventController(sseUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	studentController, err := controllers.NewStudentController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	// wait for the SSE stream to send data
	if len(studentController.GetStudentsID()) == 0 {
		for {
			if len(studentController.GetStudentsID()) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	tests := []struct {
		name           string
		server         *Server
		studentId      string
		wantStatusCode int
		wantData       bool
	}{
		// Failure
		{
			"no student controller",
			NewServer(0, eventController, nil, nil),
			"-1",
			500,
			false,
		},
		{
			"empty studentId",
			NewServer(0, eventController, nil, studentController),
			"",
			301,
			false,
		},
		{
			"invalid studentId",
			NewServer(0, eventController, nil, studentController),
			"-1",
			200,
			false,
		},

		// Success
		{
			"valid studentId",
			NewServer(0, eventController, nil, studentController),
			studentController.GetStudentsID()[0],
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/students/"+tt.studentId, nil)
			response := httptest.NewRecorder()

			Router(tt.server).ServeHTTP(response, request)
			resp := response.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Server.getStudents() gotStatusCode = %v, wantStatusCode = %v", resp.StatusCode, tt.wantStatusCode)
			}

			var gotData models.ExamResultsResponse
			_ = json.Unmarshal(body, &gotData) // intentionally ignore error

			if (len(gotData.Results) != 0) != tt.wantData {
				t.Errorf("Server.getStudents() gotData = %v, wantData = %v", len(gotData.Results) != 0, tt.wantData)
			}
		})
	}
}

func TestServer_getExams(t *testing.T) {
	eventController := controllers.NewEventController(sseUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	examController, err := controllers.NewExamController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	// wait for the SSE stream to send data
	if len(examController.GetExamNumbers()) == 0 {
		for {
			if len(examController.GetExamNumbers()) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	tests := []struct {
		name           string
		server         *Server
		wantStatusCode int
		wantData       bool
	}{
		// Failure
		{
			"no student controller",
			NewServer(0, eventController, nil, nil),
			500,
			false,
		},

		// Success
		{
			"no event controller",
			NewServer(0, nil, examController, nil),
			200,
			true,
		},
		{
			"no exam controller",
			NewServer(0, eventController, examController, nil),
			200,
			true,
		},
		{
			"valid exam controller",
			NewServer(0, eventController, examController, nil),
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/exams", nil)
			response := httptest.NewRecorder()

			Router(tt.server).ServeHTTP(response, request)
			resp := response.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Server.getExams() gotStatusCode = %v, wantStatusCode = %v", resp.StatusCode, tt.wantStatusCode)
			}

			var gotData []int
			_ = json.Unmarshal(body, &gotData) // intentionally ignore error

			if (len(gotData) != 0) != tt.wantData {
				t.Errorf("Server.getExams() gotData = %v, wantData = %v", gotData, tt.wantData)
			}
		})
	}
}

func TestServer_getExam(t *testing.T) {
	eventController := controllers.NewEventController(sseUrl)

	eventsCh, err := eventController.StartSSESubscription()
	if err != nil {
		t.Fatal(err)
	}
	defer eventController.StopSSESubscription()

	studentController, err := controllers.NewStudentController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	examController, err := controllers.NewExamController(eventsCh)
	if err != nil {
		t.Fatal(err)
	}

	// wait for the SSE stream to send data
	if len(examController.GetExamNumbers()) == 0 {
		for {
			if len(examController.GetExamNumbers()) != 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	tests := []struct {
		name           string
		server         *Server
		examNumber     int
		wantStatusCode int
		wantData       bool
	}{
		// Failure
		{
			"no exam controller",
			NewServer(0, eventController, nil, studentController),
			-1,
			500,
			false,
		},
		{
			"no student controller",
			NewServer(0, eventController, examController, nil),
			-1,
			500,
			false,
		},

		// Success
		{
			"invalid examNumber",
			NewServer(0, eventController, examController, studentController),
			-1,
			200,
			false,
		},
		{
			"valid examNumber",
			NewServer(0, eventController, examController, studentController),
			examController.GetExamNumbers()[0],
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/exams/"+strconv.Itoa(tt.examNumber), nil)
			response := httptest.NewRecorder()

			Router(tt.server).ServeHTTP(response, request)
			resp := response.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Server.getStudents() gotStatusCode = %v, wantStatusCode = %v", resp.StatusCode, tt.wantStatusCode)
			}

			var gotData models.ExamResultsResponse
			_ = json.Unmarshal(body, &gotData) // intentionally ignore error

			if (len(gotData.Results) != 0) != tt.wantData {
				t.Errorf("Server.getStudents() gotData = %v, wantData = %v", len(gotData.Results) != 0, tt.wantData)
			}
		})
	}
}
