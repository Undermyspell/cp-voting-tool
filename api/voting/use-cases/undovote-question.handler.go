package usecases

import (
	"encoding/json"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	shared_models "voting/shared/models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
	usecases_events "voting/voting/use-cases/_events"
)

func UndoVote(questionId string, userContext shared_models.UserContext) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	votes, err := undoVote(questionId, userContext)

	if err != nil {
		return err
	}

	questionUndoUpvoteMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}

	questionUndoUpVoteForUserMessage := usecases_events.QuestionUpvoted{
		Id:    questionId,
		Votes: votes,
		Voted: false,
	}

	questionForUserPaylod, errf := json.Marshal(questionUndoUpVoteForUserMessage)
	questionPayload, errj := json.Marshal(questionUndoUpvoteMessage)

	if errj != nil || errf != nil {
		return &errors.UnexpectedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "cant marshal question",
			},
		}
	}

	event := shared.Event{
		EventType: usecases_events.UNDO_UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	userevent := shared.Event{
		EventType: usecases_events.UNDO_UPVOTE_QUESTION,
		Payload:   string(questionForUserPaylod),
	}

	broker.NotifyUser(userevent, userContext)
	broker.NotifyAllButUser(event, userContext)

	return nil
}

func undoVote(id string, user shared_models.UserContext) (int, errors.VotingError) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return 0, &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	question, ok := votingStorage.GetQuestion(id)

	if !ok {
		return 0, &errors.QuestionNotFoundError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question not found",
			},
		}
	}

	answered := question.Answered
	if answered {
		return 0, &errors.QuestionAlreadyAnsweredError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "question already answered",
			},
		}
	}

	hash := user.GetHash(votingStorage.GetSecret())
	_, ok = votingStorage.GetUserVotes()[hash][id]

	if !ok {
		return 0,
			&errors.UserHasNotVotedError{
				UseCaseError: shared.UseCaseError{
					ErrMessage: "user has not voted",
				},
			}
	}

	votingStorage.UndoVote(hash, id)

	question, _ = votingStorage.GetQuestion(question.Id)

	return question.Votes, nil
}
