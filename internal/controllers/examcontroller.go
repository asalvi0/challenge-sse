package controllers

import "github.com/asalvi0/challenge-sse/internal/models"

type ExamController struct {
	exams    models.ExamsRepository
	eventsCh chan models.Event
}

func NewExamController(eventsCh chan models.Event) *ExamController {
	controller := ExamController{
		exams:    make(models.ExamsRepository, 0),
		eventsCh: eventsCh,
	}
	controller.Listen()

	return &controller
}

func (c *ExamController) Listen() {
	go func() {
		for event := range c.eventsCh {
			c.AddExam(event)
		}
	}()
}

func (c *ExamController) AddExam(event models.Event) {
	// append the exam data to the map, empty struct uses no memory
	exam := c.exams[event.Number]
	if exam == nil {
		c.exams[event.Number] = &models.ExamRecord{}
		exam = c.exams[event.Number]
	}

	// update the exam summary
	exam.Count += 1
	exam.Average = (exam.Average + event.Score) / float32(exam.Count)
}

func (c *ExamController) GetExamNumbers() (exams []int) {
	exams = make([]int, len(c.exams))

	i := 0
	for key := range c.exams {
		exams[i] = key
		i++
	}

	return exams
}

func (c *ExamController) GetExam(number int, students []models.StudentRecord) (scores []float32, average float32) {
	if c.exams[number] == nil || len(students) == 0 {
		return nil, 0
	}

	for _, student := range students {
		for i := 0; i < len(student.Exams); i++ {
			if student.Exams[i].Number == number {
				scores = append(scores, student.Exams[i].Score)
			}
		}
	}

	if len(scores) > 0 {
		return scores, c.exams[number].Average
	}
	return nil, 0
}

func (c *ExamController) GetExams() (exams []models.ExamRecord) {
	for _, exam := range c.exams {
		exams = append(exams, *exam)
	}

	return exams
}
