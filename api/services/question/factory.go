package services

import (
	"sse/internal/broker"
	"sse/internal/votingstorage"
)

func NewBrokered(broker broker.Broker, storage votingstorage.VotingStorage) QuestionService {
	return &BrokeredQuestionsService{
		QuestionSession: storage,
		Broker:          broker,
	}
}
