package usecases

import (
	"encoding/json"
	"fmt"
	"voting/dtos"
	"voting/internal/events"
	"voting/internal/models"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func Create(newQuestionDto dtos.NewQuestionDto, userContext shared_models.UserContext) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	question, err := create(newQuestionDto.Text, newQuestionDto.Anonymous, userContext)

	if err != nil {
		return err
	}

	newQuestionForUserSseMessage := events.QuestionCreated{
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

	newQuestionForAllButUserSseMessage := events.QuestionCreated{
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

	eventForUser := events.Event{
		EventType: events.NEW_QUESTION,
		Payload:   string(newQuestionForUserByteString),
	}
	eventForAllButUser := events.Event{
		EventType: events.NEW_QUESTION,
		Payload:   string(newQuestionForAllButUserByteString),
	}

	broker.NotifyUser(eventForUser, userContext)
	broker.NotifyAllButUser(eventForAllButUser, userContext)

	return nil
}

func create(text string, anonymous bool, creator shared_models.UserContext) (*models.Question, errors.VotingError) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return nil, &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	if len(text) > models.MaxLength {
		return nil, &errors.QuestionMaxLengthExceededError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: fmt.Sprintf("Question text must have a max length of %d", models.MaxLength),
			},
		}
	}

	question := votingStorage.AddQuestion(text, anonymous, creator.Name, creator.GetHash(votingStorage.GetSecret()))
	return &question, nil
}
