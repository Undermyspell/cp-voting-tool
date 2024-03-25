package usecases

import (
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
	usecases_events "voting/voting/use-cases/_events"
)

func StartSession() errors.VotingError {
	broker := shared_infra_broker.GetInstance()
	votingStorage := voting_repositories.GetInstance()
	votingStorage.Start()

	event := usecases_events.Event{
		EventType: usecases_events.START_SESSION,
		Payload:   usecases_events.PayloadEmpty,
	}

	broker.NotifyAll(event)

	return nil
}
