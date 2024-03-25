package usecases

import (
	"encoding/json"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_models "voting/voting/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
	usecases_events "voting/voting/use-cases/_events"
)

type UpdateQuestionDto struct {
	Id        string `json:"id" binding:"required"`
	Text      string `json:"text" binding:"required"`
	Anonymous bool   `json:"anonymous"`
}

func UpdateQuestion(updateQuestionDto UpdateQuestionDto, creator shared_models.UserContext) errors.VotingError {

	broker := shared_infra_broker.GetInstance()

	questionToUpdate, err := update(updateQuestionDto, creator)

	if err != nil {
		return err
	}

	questionToUpdateSseMessage := usecases_events.QuestionUpdated{
		Id:        questionToUpdate.Id,
		Text:      questionToUpdate.Text,
		Creator:   questionToUpdate.CreatorName,
		Anonymous: questionToUpdate.Anonymous,
	}

	newQuestionByteString, _ := json.Marshal(questionToUpdateSseMessage)

	event := usecases_events.Event{
		EventType: usecases_events.UPDATE_QUESTION,
		Payload:   string(newQuestionByteString),
	}

	broker.NotifyAll(event)

	return nil
}

func update(question UpdateQuestionDto, creator shared_models.UserContext) (*voting_models.Question, errors.VotingError) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return nil, &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	updatedQuestion, ok := votingStorage.GetQuestion(question.Id)
	if !ok {
		return nil, &errors.QuestionNotFoundError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question not found",
			},
		}
	}

	if updatedQuestion.CreatorHash != creator.GetHash(votingStorage.GetSecret()) {
		return nil, &errors.QuestionNotOwnedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "you do not own this question",
			},
		}
	}

	if updatedQuestion.Answered {
		return nil, &errors.QuestionAlreadyAnsweredError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question has already been answered",
			},
		}
	}

	updatedQuestion = votingStorage.UpdateQuestion(question.Id, question.Text, creator.Name, question.Anonymous)

	return &updatedQuestion, nil
}
