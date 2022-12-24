package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"githb.com/asalvi0/challenge-sse/internal/model"
	"githb.com/asalvi0/challenge-sse/internal/services"

	"github.com/goccy/go-json"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	port    int
	service *services.ExamService
}

func NewServer(port int, service *services.ExamService) *Server {
	return &Server{port, service}
}

func (s *Server) Start() {
	router := httprouter.New()

	router.GET("/students", s.getStudents)
	router.GET("/students/:id", s.getStudent)

	router.GET("/exams", s.getExams)
	router.GET("/exams/:number", s.getExam)

	router.GET("/stop-sse", s.stopSSESubscription)

	// TODO: move panic to main.go
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), router))
}

func (s *Server) stopSSESubscription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.service.StopSSESubscription()

	fmt.Fprintf(w, string("OK"))
}

func (s *Server) getStudents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	students := s.service.GetStudents()

	resp, err := json.Marshal(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) getStudent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "url param 'id' is missing")
		return
	}

	scores, average := s.service.GetStudent(id)
	student := model.StudentResponse{
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
	exams := s.service.GetExams()

	resp, err := json.Marshal(exams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, string(resp))
}

func (s *Server) getExam(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pNumber := params.ByName("number")
	if len(pNumber) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "url param 'number' is missing")
		return
	}

	number, err := strconv.Atoi(pNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	scores, average := s.service.GetExam(number)
	student := model.StudentResponse{
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
