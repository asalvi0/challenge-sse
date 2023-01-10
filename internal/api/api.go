package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asalvi0/challenge-sse/internal/controllers"
	"github.com/asalvi0/challenge-sse/internal/models"

	"github.com/goccy/go-json"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	port              int
	eventController   *controllers.EventController
	examController    *controllers.ExamController
	studentController *controllers.StudentController
}

func NewServer(port int,
	eventController *controllers.EventController,
	examController *controllers.ExamController,
	studentController *controllers.StudentController,
) *Server {
	return &Server{
		port,
		eventController,
		examController,
		studentController,
	}
}

func (s *Server) Start() {
	router := httprouter.New()

	router.GET("/students", s.getStudents)
	router.GET("/students/:id", s.getStudent)

	router.GET("/exams", s.getExams)
	router.GET("/exams/:number", s.getExam)

	router.GET("/start-sse", s.startSSESubscription)
	router.GET("/stop-sse", s.stopSSESubscription)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), router))
}

func (s *Server) getStudents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if s == nil || s.studentController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	students := s.studentController.GetStudentsID()

	resp, err := json.Marshal(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) getStudent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if s == nil || s.studentController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	id := params.ByName("id")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "url param 'id' is missing")
		return
	}

	scores, average := s.studentController.GetStudent(id)
	student := models.ExamResultsResponse{
		Results: scores,
		Average: average,
	}

	resp, err := json.Marshal(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) getExams(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if s == nil || s.examController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	exams := s.examController.GetExamNumbers()

	resp, err := json.Marshal(exams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) getExam(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if s == nil || s.examController == nil || s.studentController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	number, err := strconv.Atoi(params.ByName("number"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid parameter")
		return
	}

	students := s.studentController.GetStudents()
	scores, average := s.examController.GetExam(number, students)
	student := models.ExamResultsResponse{
		Results: scores,
		Average: average,
	}

	resp, err := json.Marshal(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) startSSESubscription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if s == nil || s.eventController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	s.eventController.StartSSESubscription()

	fmt.Fprintf(w, string("OK"))
}

func (s *Server) stopSSESubscription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if s == nil || s.eventController == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	s.eventController.StopSSESubscription()

	fmt.Fprintf(w, string("OK"))
}
