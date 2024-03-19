package votinginfra

import (
	"voting/internal/broker"
	"voting/internal/votingstorage"
)

var b broker.Broker
var vs votingstorage.VotingStorage

func GetBroker() broker.Broker {
	return b
}

func GetVotingStorage() votingstorage.VotingStorage {
	return vs
}

func InitInstances(broker broker.Broker, votingStorage votingstorage.VotingStorage) {
	b = broker
	vs = votingStorage
}
