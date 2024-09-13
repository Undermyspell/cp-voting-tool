package voting_usecases

import (
	"encoding/json"
	"fmt"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

func Answer(questionId string) error {
	broker := shared_infra_broker.GetInstance()

	err := answer(questionId)

	if err != nil {
		return err
	}

	questionMessage := usecases_events.QuestionAnswered{
		Id: questionId,
	}
	questionPayload, errj := json.Marshal(questionMessage)

	if errj != nil {
		return fmt.Errorf("%w", errors.ErrUnexpected)
	}

	event := shared.Event{
		EventType: usecases_events.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	broker.NotifyAll(event)

	return nil
}

func answer(questionId string) error {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	_, ok := votingStorage.GetQuestion(questionId)

	if !ok {
		return fmt.Errorf("%w", errors.ErrQuestionNotFound)
	}

	votingStorage.AnswerQuestion(questionId)

	return nil
}
