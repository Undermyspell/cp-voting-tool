package votingsession

import (
	"sse/internal/models"
)

type VotingSession interface {
	Start()
	Stop()
	GetQuestion(id string) (*models.Question, bool)
	SetQuestion(question *models.Question)
	DeleteQuestion(id string)
	IsRunning() bool
	GetSecret() string
	GetQuestions() map[string]*models.Question
	GetUserVotes() map[string]map[string]bool
	Vote(userHash, id string)
}
