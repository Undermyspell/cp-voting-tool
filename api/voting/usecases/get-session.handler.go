package voting_usecases

import (
	"fmt"
	"reflect"
	shared_models "voting/shared/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
)

type QuestionDto struct {
	Id        string
	Text      string
	Votes     int
	Voted     bool
	Answered  bool
	Creator   string
	Anonymous bool
	Owned     bool
}

func (question QuestionDto) GetHeaders() []string {
	questionDtoType := reflect.TypeOf(QuestionDto{})
	var headers []string
	for i := 0; i < questionDtoType.NumField(); i++ {
		field := questionDtoType.Field(i)
		headers = append(headers, field.Name)
	}
	return headers
}

func GetSession(userContext *shared_models.UserContext) ([]QuestionDto, error) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return []QuestionDto{}, fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	hash := userContext.GetHash(votingStorage.GetSecret())
	questions := []QuestionDto{}

	for _, v := range votingStorage.GetQuestions() {
		creator := v.CreatorName
		owned := v.CreatorHash == userContext.GetHash(votingStorage.GetSecret())

		questions = append(questions, QuestionDto{
			Id:        v.Id,
			Text:      v.Text,
			Votes:     v.Votes,
			Answered:  v.Answered,
			Voted:     votingStorage.GetUserVotes()[hash][v.Id],
			Anonymous: v.Anonymous,
			Creator:   creator,
			Owned:     owned,
		})
	}

	return questions, nil
}
