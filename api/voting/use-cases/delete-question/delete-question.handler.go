package usecases

import (
	"encoding/json"
	"voting/internal/events"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func Delete(questionId string, creator shared_models.UserContext) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	err := delete(questionId, creator)

	if err != nil {
		return err
	}

	questionDeletedSseMessage := events.QuestionDeleted{
		Id: questionId,
	}
	questionDeletedByteString, _ := json.Marshal(questionDeletedSseMessage)

	event := events.Event{
		EventType: events.DELETE_QUESTION,
		Payload:   string(questionDeletedByteString),
	}

	broker.NotifyAll(event)

	return nil
}

func delete(id string, creator shared_models.UserContext) errors.VotingError {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	questionToDelete, ok := votingStorage.GetQuestion(id)
	if !ok {
		return &errors.QuestionNotFoundError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question not found",
			},
		}
	}

	if questionToDelete.CreatorHash != creator.GetHash(votingStorage.GetSecret()) {
		return &errors.QuestionNotOwnedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "you do not own this question",
			},
		}
	}

	if questionToDelete.Answered {
		return &errors.QuestionAlreadyAnsweredError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question has already been answered",
			},
		}
	}

	votingStorage.DeleteQuestion(questionToDelete.Id)

	return nil
}