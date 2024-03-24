package usecases

import (
	"encoding/json"
	"voting/internal/events"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func Answer(questionId string) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	err := answer(questionId)

	if err != nil {
		return err
	}

	questionMessage := events.QuestionAnswered{
		Id: questionId,
	}
	questionPayload, errj := json.Marshal(questionMessage)

	if errj != nil {
		return &errors.UnexpectedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "cant marshal question",
			},
		}
	}

	event := events.Event{
		EventType: events.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	broker.NotifyAll(event)

	return nil
}

func answer(questionId string) errors.VotingError {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	_, ok := votingStorage.GetQuestion(questionId)

	if !ok {
		return &errors.QuestionNotFoundError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question not found",
			},
		}
	}

	votingStorage.AnswerQuestion(questionId)

	return nil
}
