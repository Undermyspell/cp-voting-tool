package usecases

import shared "voting/shared"

type VotingError interface {
	IsVotingError() bool
	Error() string
}

type UnexpectedError struct {
	shared.UseCaseError
}

type QuestionNotFoundError struct {
	shared.UseCaseError
}

type QuestionAlreadyAnsweredError struct {
	shared.UseCaseError
}

type QuestionSessionNotRunningError struct {
	shared.UseCaseError
}

type UserAlreadyVotedError struct {
	shared.UseCaseError
}

type QuestionMaxLengthExceededError struct {
	shared.UseCaseError
}

type QuestionNotOwnedError struct {
	shared.UseCaseError
}

func (err *QuestionNotFoundError) IsVotingError() bool {
	return true
}

func (err *QuestionNotFoundError) Error() string {
	return err.ErrMessage
}

func (err *QuestionAlreadyAnsweredError) IsVotingError() bool {
	return true
}

func (err *QuestionAlreadyAnsweredError) Error() string {
	return err.ErrMessage
}

func (err *QuestionSessionNotRunningError) IsVotingError() bool {
	return true
}

func (err *QuestionSessionNotRunningError) Error() string {
	return err.ErrMessage
}

func (err *UserAlreadyVotedError) IsVotingError() bool {
	return true
}

func (err *UserAlreadyVotedError) Error() string {
	return err.ErrMessage
}

func (err *UnexpectedError) IsVotingError() bool {
	return true
}

func (err *UnexpectedError) Error() string {
	return err.ErrMessage
}

func (err *QuestionMaxLengthExceededError) IsVotingError() bool {
	return true
}

func (err *QuestionMaxLengthExceededError) Error() string {
	return err.ErrMessage
}

func (err *QuestionNotOwnedError) IsVotingError() bool {
	return true
}

func (err *QuestionNotOwnedError) Error() string {
	return err.ErrMessage
}
