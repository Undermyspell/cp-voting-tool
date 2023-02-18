package services

import (
	"sse/internal/broker"
)

func NewBrokered(broker broker.Broker) QuestionService {
	return &BrokeredQuestionsService{
		Broker:    broker,
		Session:   nil,
		UserVotes: nil,
	}
}
