package votingsession

import (
	"sse/internal/broker"
	"sse/internal/helper"
	"sse/internal/models"
)

type InMemory struct {
	Broker        broker.Broker
	Questions     map[string]*models.Question
	UserVotes     *models.SafeUserVotes
	SessionSecret string
}

func (session *InMemory) Start() {
	session.UserVotes = models.NewSafeUserVotes()
	session.Questions = make(map[string]*models.Question)
	session.SessionSecret = helper.GetRandomId()
}

func (session *InMemory) Stop() {
	session.UserVotes = nil
	session.Questions = nil
	session.SessionSecret = ""
}

func (session *InMemory) GetQuestion(id string) (*models.Question, bool) {
	question, ok := session.Questions[id]
	return question, ok
}

func (session *InMemory) IsRunning() bool {
	return session.Questions != nil
}

func (session *InMemory) GetSecret() string {
	return session.SessionSecret
}

func (session *InMemory) GetQuestions() map[string]*models.Question {
	return session.Questions
}

func (session *InMemory) GetUserVotes() map[string]map[string]bool {
	return session.UserVotes.Value()
}

func (session *InMemory) SetQuestion(question *models.Question) {
	session.Questions[question.Id] = question
}

func (session *InMemory) DeleteQuestion(id string) {
	delete(session.Questions, id)
}

func (session *InMemory) Vote(userHash, id string) {
	session.UserVotes.SetUserVote(userHash, id)
	session.Questions[id].Votes.Increment()
}
