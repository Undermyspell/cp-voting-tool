package models

import "sse/internal/helper"

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
		Id:        helper.GetRandomId(),
		Text:      text,
		Votes:     0,
		Answered:  false,
		Creator:   creator,
		Anonymous: anonymous,
	}
}
