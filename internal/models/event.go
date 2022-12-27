package models

import (
	"errors"
	"log"

	"github.com/goccy/go-json"
)

type (
	Event struct {
		StudentId string `json:"studentId"`
		ExamResult
	}
)

func ParseEvent(data []byte) (event Event, err error) {
	if !json.Valid(data) {
		log.Println(err)
		return Event{}, errors.New("invalid json data")
	}

	err = json.Unmarshal(data, &event)
	if err != nil {
		log.Println(err)
		return Event{}, errors.New("failed parsing data")
	}

	return event, nil
}
