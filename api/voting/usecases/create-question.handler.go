package voting_usecases

import (
	"encoding/json"
	"fmt"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_models "voting/voting/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

type NewQuestionDto struct {
	Text      string `json:"text" binding:"required"`
	Anonymous bool   `json:"anonymous"`
}

func Create(newQuestionDto NewQuestionDto, userContext shared_models.UserContext) error {
	broker := shared_infra_broker.GetInstance()

	question, err := create(newQuestionDto.Text, newQuestionDto.Anonymous, userContext)

	if err != nil {
		return err
	}

	newQuestionForUserSseMessage := usecases_events.QuestionCreated{
		Id:        question.Id,
		Text:      question.Text,
		Creator:   question.CreatorName,
		Answered:  question.Answered,
		Votes:     1,
		Voted:     true,
		Anonymous: question.Anonymous,
		Owned:     true,
	}

	creatorForAllButUser := ""

	if !question.Anonymous {
		creatorForAllButUser = question.CreatorName
	}

	newQuestionForAllButUserSseMessage := usecases_events.QuestionCreated{
		Id:        question.Id,
		Text:      question.Text,
		Creator:   creatorForAllButUser,
		Answered:  question.Answered,
		Votes:     1,
		Anonymous: question.Anonymous,
		Owned:     false,
	}

	newQuestionForUserByteString, errf := json.Marshal(newQuestionForUserSseMessage)
	newQuestionForAllButUserByteString, errj := json.Marshal(newQuestionForAllButUserSseMessage)

	if errj != nil || errf != nil {
		return fmt.Errorf("%w", errors.ErrUnexpected)
	}

	eventForUser := shared.Event{
		EventType: usecases_events.NEW_QUESTION,
		Payload:   string(newQuestionForUserByteString),
	}
	eventForAllButUser := shared.Event{
		EventType: usecases_events.NEW_QUESTION,
		Payload:   string(newQuestionForAllButUserByteString),
	}

	broker.NotifyUser(eventForUser, userContext)
	broker.NotifyAllButUser(eventForAllButUser, userContext)

	return nil
}

func create(text string, anonymous bool, creator shared_models.UserContext) (*voting_models.Question, error) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return nil, fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	if len(text) > voting_models.MaxLength {
		return nil, fmt.Errorf("%w", errors.ErrQuestionMaxLengthExceeded)
	}

	question := votingStorage.AddQuestion(text, anonymous, creator.Name, creator.GetHash(votingStorage.GetSecret()))
	return &question, nil
}
