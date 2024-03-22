package services

import (
	"voting/internal/votingstorage"
	shared_infra_broker "voting/shared/infra/broker"
)

func NewBrokered(broker shared_infra_broker.Broker, storage votingstorage.VotingStorage) QuestionService {
	return &BrokeredQuestionsService{
		QuestionSession: storage,
		Broker:          broker,
	}
}
