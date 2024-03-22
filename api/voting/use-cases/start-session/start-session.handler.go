package usecases

import (
	"voting/internal/events"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func StartSession() errors.VotingError {
	broker := shared_infra_broker.GetInstance()
	votingStorage := voting_repositories.GetInstance()
	votingStorage.Start()

	event := events.Event{
		EventType: events.START_SESSION,
		Payload:   events.PayloadEmpty,
	}

	broker.NotifyAll(event)

	return nil
}
