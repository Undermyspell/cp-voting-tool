package usecases

import (
	"voting/internal/events"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func StopSession() errors.VotingError {
	broker := shared_infra_broker.GetInstance()
	votingStorage := voting_repositories.GetInstance()
	votingStorage.Stop()

	event := events.Event{
		EventType: events.STOP_SESSION,
		Payload:   events.PayloadEmpty,
	}

	broker.NotifyAll(event)

	return nil
}
