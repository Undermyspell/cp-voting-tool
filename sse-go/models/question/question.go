package question

import "sse/internal/helper"

type Question struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Votes    int
	Answered bool
}

func New(text string) Question {
	return Question{
		Id:       helper.GetRandomId(),
		Text:     text,
		Votes:    0,
		Answered: false,
	}
}
