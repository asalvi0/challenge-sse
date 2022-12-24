package model

import (
	"errors"
	"log"

	"github.com/goccy/go-json"
)

type (
	Message struct {
		StudentId string `json:"studentId"`
		ExamResult
	}

	ExamResult struct {
		Number int     `json:"exam"`
		Score  float32 `json:"score"`
	}

	ExamRecord struct {
		Count   int     `json:"count"`
		Average float32 `json:"average"`
	}

	StudentRecord struct {
		Average float32      `json:"average"`
		Exams   []ExamResult `json:"exams"`
	}

	StudentResponse struct {
		Results []float32 `json:"results"`
		Average float32   `json:"average"`
	}

	StudentsRepository map[string]*StudentRecord
	ExamsRepository    map[int]*ExamRecord
)

func ParseExam(msg []byte) (exam Message, err error) {
	if !json.Valid(msg) {
		log.Println(err)
		return Message{}, errors.New("invalid json data")
	}

	err = json.Unmarshal(msg, &exam)
	if err != nil {
		log.Println(err)
		return Message{}, errors.New("failed parsing data")
	}

	return exam, nil
}
