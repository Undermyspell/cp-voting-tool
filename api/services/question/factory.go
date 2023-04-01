package services

import (
	"sse/internal/broker"
	"sse/internal/votingsession"
)

func NewBrokered(broker broker.Broker) QuestionService {
	return &BrokeredQuestionsService{
		QuestionSession: &votingsession.InMemory{},
		Broker:          broker,
	}
}
