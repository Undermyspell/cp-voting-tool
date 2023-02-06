package services

import (
	"sse/internal/broker"
	"sse/models"
)

func NewBrokered(broker broker.Broker) QuestionService {
	return &BrokeredQuestionsService{
		Broker:  broker,
		Session: make(map[string]models.Question),
	}
}
