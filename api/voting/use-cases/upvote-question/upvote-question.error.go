package usecases

type UpvoteQuestionError interface {
	IsUpvoteQuestionError() bool
	Error() string
}

type UseCaseError struct {
	ErrMessage string
}

type UnexpectedError struct {
	UseCaseError
}

type QuestionNotFoundError struct {
	UseCaseError
}

type QuestionAlreadyAnsweredError struct {
	UseCaseError
}

type QuestionSessionNotRunningError struct {
	UseCaseError
}

type UserAlreadyVotedError struct {
	UseCaseError
}

func (err *QuestionNotFoundError) IsUpvoteQuestionError() bool {
	return true
}

func (err *QuestionNotFoundError) Error() string {
	return err.ErrMessage
}

func (err *QuestionAlreadyAnsweredError) IsUpvoteQuestionError() bool {
	return true
}

func (err *QuestionAlreadyAnsweredError) Error() string {
	return err.ErrMessage
}

func (err *QuestionSessionNotRunningError) IsUpvoteQuestionError() bool {
	return true
}

func (err *QuestionSessionNotRunningError) Error() string {
	return err.ErrMessage
}

func (err *UserAlreadyVotedError) IsUpvoteQuestionError() bool {
	return true
}

func (err *UserAlreadyVotedError) Error() string {
	return err.ErrMessage
}

func (err *UnexpectedError) IsUpvoteQuestionError() bool {
	return true
}

func (err *UnexpectedError) Error() string {
	return err.ErrMessage
}
