package usecases

import (
	"encoding/json"
	"voting/internal/events"
	"voting/internal/models"
	votinginfra "voting/voting/infra"
)

func Upvote(questionId string, userContext models.UserContext) UpvoteQuestionError {
	broker := votinginfra.GetBroker()

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
		return &UnexpectedError{
			UseCaseError: UseCaseError{
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

func upVote(id string, user models.UserContext) (int, UpvoteQuestionError) {
	votingStorage := votinginfra.GetVotingStorage()

	if !votingStorage.IsRunning() {
		return 0, &QuestionSessionNotRunningError{
			UseCaseError: UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	question, ok := votingStorage.GetQuestion(id)

	if !ok {
		return 0, &QuestionNotFoundError{
			UseCaseError: UseCaseError{
				ErrMessage: "question not found",
			},
		}
	}

	answered := question.Answered
	if answered {
		return 0, &QuestionAlreadyAnsweredError{
			UseCaseError: UseCaseError{
				ErrMessage: "question already answered",
			},
		}
	}

	hash := user.GetHash(votingStorage.GetSecret())
	_, ok = votingStorage.GetUserVotes()[hash][id]

	if ok {
		return 0,
			&UserAlreadyVotedError{
				UseCaseError: UseCaseError{
					ErrMessage: "user already voted",
				},
			}
	}

	votingStorage.Vote(hash, id)

	question, _ = votingStorage.GetQuestion(question.Id)

	return question.Votes, nil
}
