package voting_usecases

import (
	"encoding/json"
	"fmt"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

func Delete(questionId string, creator shared_models.UserContext) error {
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

func delete(id string, creator shared_models.UserContext) error {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	questionToDelete, ok := votingStorage.GetQuestion(id)
	if !ok {
		return fmt.Errorf("%w", errors.ErrQuestionNotFound)
	}

	if questionToDelete.CreatorHash != creator.GetHash(votingStorage.GetSecret()) {
		return fmt.Errorf("%w", errors.ErrQuestionNotOwned)
	}

	if questionToDelete.Answered {
		return fmt.Errorf("%w", errors.ErrQuestionAlreadyAnswered)
	}

	votingStorage.DeleteQuestion(questionToDelete.Id)

	return nil
}
