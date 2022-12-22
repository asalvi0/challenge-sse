package model

type Data struct {
	Exam      int     `json:"exam,omitempty"`
	StudentId string  `json:"studentId,omitempty"`
	Score     float32 `json:"score,omitempty"`
}
