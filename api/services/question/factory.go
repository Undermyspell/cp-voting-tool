package services

import (
	"voting/internal/broker"
	"voting/internal/votingstorage"
)

func NewBrokered(broker broker.Broker, storage votingstorage.VotingStorage) QuestionService {
	return &BrokeredQuestionsService{
		QuestionSession: storage,
		Broker:          broker,
	}
}
