package usecases

import (
	"encoding/json"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
	usecases_events "voting/voting/use-cases/_events"
)

func Delete(questionId string, creator shared_models.UserContext) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	err := delete(questionId, creator)

	if err != nil {
		return err
	}

	questionDeletedSseMessage := usecases_events.QuestionDeleted{
		Id: questionId,
	}
	questionDeletedByteString, _ := json.Marshal(questionDeletedSseMessage)

	event := shared.Event{
		EventType: usecases_events.DELETE_QUESTION,
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
