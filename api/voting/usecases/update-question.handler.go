package voting_usecases

import (
	"encoding/json"
	"fmt"
	"voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_models "voting/voting/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

type UpdateQuestionDto struct {
	Id        string `json:"id" binding:"required"`
	Text      string `json:"text" binding:"required"`
	Anonymous bool   `json:"anonymous"`
}

func UpdateQuestion(updateQuestionDto UpdateQuestionDto, creator shared_models.UserContext) error {

	broker := shared_infra_broker.GetInstance()

	questionToUpdate, err := update(updateQuestionDto, creator)

	if err != nil {
		return err
	}

	questionToUpdateSseMessage := usecases_events.QuestionUpdated{
		Id:        questionToUpdate.Id,
		Text:      questionToUpdate.Text,
		Creator:   questionToUpdate.CreatorName,
		Anonymous: questionToUpdate.Anonymous,
	}

	newQuestionByteString, _ := json.Marshal(questionToUpdateSseMessage)

	event := shared.Event{
		EventType: usecases_events.UPDATE_QUESTION,
		Payload:   string(newQuestionByteString),
	}

	broker.NotifyAll(event)

	return nil
}

func update(question UpdateQuestionDto, creator shared_models.UserContext) (*voting_models.Question, error) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return nil, fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	updatedQuestion, ok := votingStorage.GetQuestion(question.Id)
	if !ok {
		return nil, fmt.Errorf("%w", errors.ErrQuestionNotFound)
	}

	if updatedQuestion.CreatorHash != creator.GetHash(votingStorage.GetSecret()) {
		return nil, fmt.Errorf("%w", errors.ErrQuestionNotOwned)
	}

	if updatedQuestion.Answered {
		return nil, fmt.Errorf("%w", errors.ErrQuestionAlreadyAnswered)
	}

	updatedQuestion = votingStorage.UpdateQuestion(question.Id, question.Text, creator.Name, question.Anonymous)

	return &updatedQuestion, nil
}
