package models

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

type Question struct {
	Id          string `json:"id"`
	Text        string `json:"text"`
	Votes       SafeCounter
	Answered    bool
	Voted       bool
	CreatorHash string
	CreatorName string
	Anonymous   bool
}

func NewQuestion(text string, anonymous bool, creatorName, creatorHash string) *Question {
	if anonymous {
		creatorName = ""
	}
	return &Question{
		Id:          ulid.Make().String(),
		Text:        text,
		Votes:       SafeCounter{mu: sync.Mutex{}},
		Answered:    false,
		Voted:       false,
		CreatorHash: creatorHash,
		CreatorName: creatorName,
		Anonymous:   anonymous,
	}
}
