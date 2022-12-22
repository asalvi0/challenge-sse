package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Port int
}

func NewServer(port int) *Server {
	return &Server{port}
}

func (s *Server) Start() {
	router := httprouter.New()

	router.GET("/students", students)
	router.GET("/students/:id", student)

	router.GET("/exams", exams)
	router.GET("/exams/:id", exam)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), router))
}

func students(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Students")
}

func student(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if len(id) == 0 {
		fmt.Fprintf(w, "url param 'id' is missing")
		return
	}

	fmt.Fprintf(w, "Student: %s", id)
}

func exams(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Exams")
}

func exam(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	number := params.ByName("number")
	if len(number) == 0 {
		fmt.Fprintf(w, "url param 'number' is missing")
		return
	}

	fmt.Fprintf(w, "Exam: %s", number)
}
