package voting_usecases

import (
	"encoding/json"
	"fmt"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/usecases/_errors"
	usecases_events "voting/voting/usecases/_events"
)

func Upvote(questionId string, userContext shared_models.UserContext) error {
	broker := shared_infra_broker.GetInstance()

	votes, err := upVote(questionId, userContext)

	if err != nil {
		return err
	}

	questionUpvoteMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}

	questionUpVoteForUserMessage := usecases_events.QuestionUpvoted{
		Id:    questionId,
		Votes: votes,
		Voted: true,
	}

	questionForUserPaylod, errf := json.Marshal(questionUpVoteForUserMessage)
	questionPayload, errj := json.Marshal(questionUpvoteMessage)

	if errj != nil || errf != nil {
		return fmt.Errorf("%w", errors.ErrUnexpected)
	}

	event := shared.Event{
		EventType: usecases_events.UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	userevent := shared.Event{
		EventType: usecases_events.UPVOTE_QUESTION,
		Payload:   string(questionForUserPaylod),
	}

	broker.NotifyUser(userevent, userContext)
	broker.NotifyAllButUser(event, userContext)

	return nil
}

func upVote(id string, user shared_models.UserContext) (int, error) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return 0, fmt.Errorf("%w", errors.ErrQuestionSessionNotRunning)
	}

	question, ok := votingStorage.GetQuestion(id)

	if !ok {
		return 0, fmt.Errorf("%w", errors.ErrQuestionNotFound)
	}

	answered := question.Answered
	if answered {
		return 0, fmt.Errorf("%w", errors.ErrQuestionAlreadyAnswered)
	}

	hash := user.GetHash(votingStorage.GetSecret())
	_, ok = votingStorage.GetUserVotes()[hash][id]

	if ok {
		return 0, fmt.Errorf("%w", errors.ErrUserAlreadyVoted)
	}

	votingStorage.Vote(hash, id)

	question, _ = votingStorage.GetQuestion(question.Id)

	return question.Votes, nil
}
