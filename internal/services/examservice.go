package services

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"githb.com/asalvi0/challenge-sse/internal/model"

	"github.com/r3labs/sse/v2"
)

type ExamService struct {
	// students can take the same exam multiple times
	students model.StudentsRepository

	// deduplicated exam numbers
	exams model.ExamsRepository

	// SSE
	sseClient *sse.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewExamService() *ExamService {
	service := ExamService{
		students: make(model.StudentsRepository, 0),
		exams:    make(model.ExamsRepository, 0),
	}
	service.ctx, service.ctxCancel = context.WithCancel(context.Background())

	return &service
}

func validateSSEUrl(sseUrl string) (err error) {
	// validate url format
	parsedUrl, err := url.Parse(sseUrl)
	if (err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") || parsedUrl.Host == "") || len(sseUrl) == 0 {
		return errors.New("invalid SSE url")
	}

	// check if url is reachable
	// TODO: implement a proper validation method
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}

	resp, err := client.Get(sseUrl)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("unreachable SSE url")
	}

	return nil
}

func (s *ExamService) SubscribeSSE(url string) (err error) {
	err = validateSSEUrl(url)
	if err != nil {
		return err
	}

	s.sseClient = sse.NewClient(url)
	if s.sseClient == nil {
		return errors.New("failed to create SSE client")
	}

	go func() {
		s.sseClient.SubscribeRawWithContext(s.ctx, func(msg *sse.Event) {
			err := s.ProcessMessage(msg.Data)
			if err != nil {
				log.Println(err)
			}
		})
	}()

	return nil
}

func (s *ExamService) StopSSESubscription() {
	s.ctxCancel()
}

func (s *ExamService) ProcessMessage(message []byte) (err error) {
	record, err := model.ParseExam(message)
	if err != nil {
		return err
	}

	if len(record.StudentId) == 0 || record.Number <= 0 {
		return errors.New("invalid message data")
	}

	// append the exam data for the student
	student := s.students[record.StudentId]
	if student == nil {
		s.students[record.StudentId] = &model.StudentRecord{}
		student = s.students[record.StudentId]
	}

	// update the student record
	student.Exams = append(student.Exams, record.ExamResult)
	student.Average += (student.Average + record.Score) / float32(len(student.Exams))

	// append the exam data to the map, empty struct uses no memory
	exam := s.exams[record.Number]
	if exam == nil {
		s.exams[record.Number] = &model.ExamRecord{}
		exam = s.exams[record.Number]
	}

	// update the exam summary
	exam.Count += 1
	exam.Average = (exam.Average + record.Score) / float32(exam.Count)

	return nil
}

func (s *ExamService) GetStudents() (students []string) {
	students = make([]string, len(s.students))

	i := 0
	for key := range s.students {
		students[i] = key
		i++
	}

	return students
}

func (s *ExamService) GetStudent(id string) (scores []float32, average float32) {
	if s.students[id] == nil || len(s.students[id].Exams) == 0 {
		return nil, 0
	}

	for i := 0; i < len(s.students[id].Exams); i++ {
		scores = append(scores, s.students[id].Exams[i].Score)
		average += scores[i]
	}
	average = average / float32(len(scores))

	return scores, average
}

func (s *ExamService) GetExams() (exams []int) {
	exams = make([]int, len(s.exams))

	i := 0
	for key := range s.exams {
		exams[i] = key
		i++
	}

	return exams
}

func (s *ExamService) GetExam(number int) (scores []float32, average float32) {
	if s.exams[number] == nil || len(s.students) == 0 {
		return nil, 0
	}

	for _, student := range s.students {
		for i := 0; i < len(student.Exams); i++ {
			if student.Exams[i].Number == number {
				scores = append(scores, student.Exams[i].Score)
			}
		}
	}

	if len(scores) > 0 {
		return scores, s.exams[number].Average
	}
	return nil, 0
}
