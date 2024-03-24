package usecases

import (
	"voting/dtos"
	"voting/shared"
	"voting/shared/shared_models"
	voting_repositories "voting/voting/repositories"
	errors "voting/voting/use-cases/_errors"
)

func GetSession(userContext *shared_models.UserContext) ([]dtos.QuestionDto, errors.VotingError) {
	votingStorage := voting_repositories.GetInstance()

	if !votingStorage.IsRunning() {
		return []dtos.QuestionDto{}, &errors.QuestionSessionNotRunningError{
			UseCaseError: shared.UseCaseError{
				ErrMessage: "no questions session currently running",
			},
		}
	}

	hash := userContext.GetHash(votingStorage.GetSecret())
	questions := []dtos.QuestionDto{}

	for _, v := range votingStorage.GetQuestions() {
		creator := v.CreatorName
		owned := v.CreatorHash == userContext.GetHash(votingStorage.GetSecret())

		questions = append(questions, dtos.QuestionDto{
			Id:        v.Id,
			Text:      v.Text,
			Votes:     v.Votes,
			Answered:  v.Answered,
			Voted:     votingStorage.GetUserVotes()[hash][v.Id],
			Anonymous: v.Anonymous,
			Creator:   creator,
			Owned:     owned,
		})
	}

	return questions, nil
}
