package models

import (
	"github.com/oklog/ulid/v2"
)

type Question struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Votes     int
	Answered  bool
	Creator   UserContext
	Anonymous bool
}

func NewQuestion(text string, anonymous bool, creator UserContext) Question {
	return Question{
		Id:        ulid.Make().String(),
		Text:      text,
		Votes:     0,
		Answered:  false,
		Creator:   creator,
		Anonymous: anonymous,
	}
}
