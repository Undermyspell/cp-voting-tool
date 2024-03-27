package voting_usecases

import (
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

func StartSession() errors.VotingError {
	broker := shared_infra_broker.GetInstance()
	votingStorage := voting_repositories.GetInstance()
	votingStorage.Start()

	event := shared.Event{
		EventType: usecases_events.START_SESSION,
		Payload:   usecases_events.PayloadEmpty,
	}

	broker.NotifyAll(event)

	return nil
}
