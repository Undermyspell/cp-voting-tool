package usecases

import (
	"encoding/json"
	"voting/internal/events"
	shared "voting/shared"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func Upvote(questionId string, userContext shared_models.UserContext) errors.VotingError {
	broker := shared_infra_broker.GetInstance()

	votes, err := upVote(questionId, userContext)

	if err != nil {
		return err
	}

	questionUpvoteMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}

	questionUpVoteForUserMessage := events.QuestionUpvoted{
		Id:    questionId,
		Votes: votes,
		Voted: true,
	}

	questionForUserPaylod, errf := json.Marshal(questionUpVoteForUserMessage)
	questionPayload, errj := json.Marshal(questionUpvoteMessage)

	if errj != nil || errf != nil {
		return &errors.UnexpectedError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "cant marshal question",
			},
		}
	}

	event := events.Event{
		EventType: events.UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	userevent := events.Event{
		EventType: events.UPVOTE_QUESTION,
		Payload:   string(questionForUserPaylod),
	}

	broker.NotifyUser(userevent, userContext)
	broker.NotifyAllButUser(event, userContext)

	return nil
}

func upVote(id string, user shared_models.UserContext) (int, errors.VotingError) {
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

	if ok {
		return 0,
			&errors.UserAlreadyVotedError{
				UseCaseError: shared.UseCaseError{
					ErrMessage: "user already voted",
				},
			}
	}

	votingStorage.Vote(hash, id)

	question, _ = votingStorage.GetQuestion(question.Id)

	return question.Votes, nil
}
