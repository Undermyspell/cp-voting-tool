package voting_repositories

import "voting/internal/votingstorage"

var vs votingstorage.VotingStorage

func GetInstance() votingstorage.VotingStorage {
	return vs
}

func InitInstances(votingStorage votingstorage.VotingStorage) {
	vs = votingStorage
}
