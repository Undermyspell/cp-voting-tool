package models

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

type Question struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Votes     SafeCounter
	Voted     bool
	Answered  bool
	Creator   UserContext
	Anonymous bool
}

func NewQuestion(text string, anonymous bool, creator UserContext) *Question {
	return &Question{
		Id:        ulid.Make().String(),
		Text:      text,
		Votes:     SafeCounter{mu: sync.Mutex{}},
		Answered:  false,
		Voted:     false,
		Creator:   creator,
		Anonymous: anonymous,
	}
}
