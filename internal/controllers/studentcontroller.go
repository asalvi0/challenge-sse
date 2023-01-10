package controllers

import (
	"errors"

	"github.com/asalvi0/challenge-sse/internal/models"
)

type StudentController struct {
	students models.StudentsRepository
	eventsCh <-chan models.Event
}

func NewStudentController(eventsCh <-chan models.Event) (*StudentController, error) {
	if eventsCh == nil {
		return nil, errors.New("nil channel is not allowed")
	}

	controller := StudentController{
		students: make(models.StudentsRepository, 0),
		eventsCh: eventsCh,
	}
	controller.listen()

	return &controller, nil
}

func (c *StudentController) listen() {
	go func() {
		for event := range c.eventsCh {
			c.AddStudent(event)
		}
	}()
}

func (c *StudentController) AddStudent(event models.Event) error {
	if len(event.StudentId) <= 0 {
		return errors.New("invalid event: studentId is required")
	}
	if event.Number == 0 {
		return errors.New("invalid event: exam number is required")
	}
	if event.Score <= 0 {
		return errors.New("invalid event: exam score is required")
	}

	// append the exam data for the student
	student := c.students[event.StudentId]
	if student == nil {
		c.students[event.StudentId] = &models.StudentRecord{}
		student = c.students[event.StudentId]
		student.Id = event.StudentId
	}

	// update the student record
	student.Exams = append(student.Exams, event.ExamResult)
	student.Average += (student.Average + event.Score) / float32(len(student.Exams))

	return nil
}

func (c *StudentController) GetStudentsID() (students []string) {
	students = make([]string, len(c.students))

	i := 0
	for key := range c.students {
		students[i] = key
		i++
	}

	return students
}

func (c *StudentController) GetStudent(id string) (scores []float32, average float32) {
	if c.students[id] == nil || len(c.students[id].Exams) == 0 {
		return nil, 0
	}

	for i := 0; i < len(c.students[id].Exams); i++ {
		scores = append(scores, c.students[id].Exams[i].Score)
		average += scores[i]
	}
	average = average / float32(len(scores))

	return scores, average
}

func (c *StudentController) GetStudents() (students []models.StudentRecord) {
	if len(c.students) == 0 {
		return students
	}

	for _, student := range c.students {
		students = append(students, *student)
	}

	return students
}
