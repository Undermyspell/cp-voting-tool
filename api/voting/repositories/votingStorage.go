package voting_repositories

import voting_models "voting/voting/models"

var vs VotingStorage

func GetInstance() VotingStorage {
	return vs
}

func InitInstances(votingStorage VotingStorage) {
	vs = votingStorage
}

type VotingStorage interface {
	Start()
	Stop()
	GetQuestion(id string) (voting_models.Question, bool)
	AddQuestion(text string, anonymous bool, creatorName, creatorHash string) voting_models.Question
	UpdateQuestion(id, text, creatorName string, anonymous bool) voting_models.Question
	DeleteQuestion(id string)
	AnswerQuestion(id string)
	IsRunning() bool
	GetSecret() string
	GetQuestions() map[string]voting_models.Question
	GetUserVotes() map[string]map[string]bool
	Vote(userHash, id string)
	UndoVote(userHash, id string)
}
