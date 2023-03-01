package models

import "sync"

type SafeUserVotes struct {
	mu    sync.Mutex
	value map[string]map[string]bool
}

func (userVotes *SafeUserVotes) SetUserVote(userHash, questionId string) {
	userVotes.mu.Lock()
	defer userVotes.mu.Unlock()
	_, ok := userVotes.value[userHash]

	if !ok {
		userVotes.value[userHash] = make(map[string]bool)
	}

	userVotes.value[userHash][questionId] = true
}

func (userVotes *SafeUserVotes) Value() map[string]map[string]bool {
	userVotes.mu.Lock()
	defer userVotes.mu.Unlock()
	return userVotes.value
}

func NewSafeUserVotes() *SafeUserVotes {
	return &SafeUserVotes{
		value: make(map[string]map[string]bool),
	}
}
