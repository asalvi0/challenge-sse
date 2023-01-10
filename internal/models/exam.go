package models

type (
	ExamResult struct {
		Number int     `json:"exam"`
		Score  float32 `json:"score"`
	}

	ExamResultsResponse struct {
		Results []float32 `json:"results"`
		Average float32   `json:"average"`
	}

	ExamRecord struct {
		Count   int     `json:"count"`
		Average float32 `json:"average"`
	}

	ExamsRepository map[int]*ExamRecord
)
