package usecases

import (
	"encoding/json"
	"fmt"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_models "voting/voting/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
	usecases_events "voting/voting/use-cases/_events"
)

type NewQuestionDto struct {
	Text      string `json:"text" binding:"required"`
	Anonymous bool   `json:"anonymous"`
}

func Create(newQuestionDto NewQuestionDto, userContext shared_models.UserContext) errors.VotingError {
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
		return &errors.UnexpectedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "cant marshal question",
			},
		}
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

func create(text string, anonymous bool, creator shared_models.UserContext) (*voting_models.Question, errors.VotingError) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return nil, &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	if len(text) > voting_models.MaxLength {
		return nil, &errors.QuestionMaxLengthExceededError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: fmt.Sprintf("Question text must have a max length of %d", voting_models.MaxLength),
			},
		}
	}

	question := votingStorage.AddQuestion(text, anonymous, creator.Name, creator.GetHash(votingStorage.GetSecret()))
	return &question, nil
}
