package votingstorage

import (
	"sse/internal/helper"
	"sse/internal/models"

	"github.com/oklog/ulid/v2"
)

type questionInMemory struct {
	Id          string `json:"id"`
	Text        string `json:"text"`
	Votes       models.SafeCounter
	Answered    bool
	CreatorHash string
	CreatorName string
	Anonymous   bool
}

func newQuestion(text string, anonymous bool, creatorName, creatorHash string) *questionInMemory {
	if anonymous {
		creatorName = ""
	}
	return &questionInMemory{
		Id:          ulid.Make().String(),
		Text:        text,
		Votes:       models.SafeCounter{},
		Answered:    false,
		CreatorHash: creatorHash,
		CreatorName: creatorName,
		Anonymous:   anonymous,
	}
}

type InMemory struct {
	Questions     map[string]*questionInMemory
	UserVotes     *models.SafeUserVotes
	SessionSecret string
}

func (session *InMemory) Start() {
	session.UserVotes = models.NewSafeUserVotes()
	session.Questions = make(map[string]*questionInMemory)
	session.SessionSecret = helper.GetRandomId(30)
}

func (session *InMemory) Stop() {
	session.UserVotes = nil
	session.Questions = nil
	session.SessionSecret = ""
}

func (session *InMemory) GetQuestion(id string) (models.Question, bool) {
	questionInMemory, ok := session.Questions[id]

	question := models.Question{}
	if ok {
		question = questionFromInMemoryQuestion(questionInMemory)
	}

	return question, ok
}

func (session *InMemory) IsRunning() bool {
	return session.Questions != nil
}

func (session *InMemory) GetSecret() string {
	return session.SessionSecret
}

func (session *InMemory) GetQuestions() map[string]models.Question {
	questions := make(map[string]models.Question)

	for id, question := range session.Questions {
		questions[id] = questionFromInMemoryQuestion(question)
	}

	return questions
}

func (session *InMemory) GetUserVotes() map[string]map[string]bool {
	return session.UserVotes.Value()
}

func (session *InMemory) AddQuestion(text string, anonymous bool, creatorName, creatorHash string) models.Question {
	question := newQuestion(text, anonymous, creatorName, creatorHash)
	session.Questions[question.Id] = question
	session.Vote(creatorHash, question.Id)
	return questionFromInMemoryQuestion(question)
}

func (session *InMemory) UpdateQuestion(id, text, creatorName string, anonymous bool) models.Question {
	questionToUpdate := session.Questions[id]
	questionToUpdate.Text = text
	questionToUpdate.Anonymous = anonymous
	questionToUpdate.CreatorName = ""
	if !questionToUpdate.Anonymous {
		questionToUpdate.CreatorName = creatorName
	}
	return questionFromInMemoryQuestion(questionToUpdate)
}

func (session *InMemory) AnswerQuestion(id string) {
	questionToAnswer := session.Questions[id]
	questionToAnswer.Answered = true
}

func (session *InMemory) DeleteQuestion(id string) {
	delete(session.Questions, id)
}

func (session *InMemory) Vote(userHash, id string) {
	session.UserVotes.SetUserVote(userHash, id)
	session.Questions[id].Votes.Increment()
}

func (session *InMemory) UndoVote(userHash, id string) {
	session.UserVotes.RemoveUserVote(userHash, id)
	session.Questions[id].Votes.Decrement()
}

func questionFromInMemoryQuestion(question *questionInMemory) models.Question {
	return models.NewQuestion(
		question.Id,
		question.Text,
		question.Votes.Value(),
		question.Answered,
		false,
		question.Anonymous,
		question.CreatorName,
		question.CreatorHash)
}
