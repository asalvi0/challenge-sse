package models

type (
	StudentRecord struct {
		Id      string
		Average float32      `json:"average"`
		Exams   []ExamResult `json:"exams"`
	}

	StudentsRepository map[string]*StudentRecord
)
