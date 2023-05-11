package votingstorage

import (
	"sse/internal/models"
)

type VotingStorage interface {
	Start()
	Stop()
	GetQuestion(id string) (models.Question, bool)
	AddQuestion(text string, anonymous bool, creatorName, creatorHash string) models.Question
	UpdateQuestion(id, text, creatorName string, anonymous bool) models.Question
	DeleteQuestion(id string)
	AnswerQuestion(id string)
	IsRunning() bool
	GetSecret() string
	GetQuestions() map[string]models.Question
	GetUserVotes() map[string]map[string]bool
	Vote(userHash, id string)
	UndoVote(userHash, id string)
}
