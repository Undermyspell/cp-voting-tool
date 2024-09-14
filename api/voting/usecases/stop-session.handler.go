package voting_usecases

import (
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	voting_repositories "voting/voting/repositories"
	usecases_events "voting/voting/usecases/_events"
)

func StopSession() error {
	broker := shared_infra_broker.GetInstance()
	votingStorage := voting_repositories.GetInstance()
	votingStorage.Stop()

	event := shared.Event{
		EventType: usecases_events.STOP_SESSION,
		Payload:   usecases_events.PayloadEmpty,
	}

	broker.NotifyAll(event)

	return nil
}
