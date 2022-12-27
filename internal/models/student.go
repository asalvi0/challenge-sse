package models

type (
	StudentRecord struct {
		Average float32      `json:"average"`
		Exams   []ExamResult `json:"exams"`
	}

	StudentResponse struct {
		Results []float32 `json:"results"`
		Average float32   `json:"average"`
	}

	StudentsRepository map[string]*StudentRecord
)
