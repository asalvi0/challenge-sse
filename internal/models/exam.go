package models

type (
	ExamResult struct {
		Number int     `json:"exam"`
		Score  float32 `json:"score"`
	}

	ExamRecord struct {
		Count   int     `json:"count"`
		Average float32 `json:"average"`
	}

	ExamsRepository map[int]*ExamRecord
)
