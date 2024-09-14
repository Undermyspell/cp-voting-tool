package voting_errors

import (
	"errors"
)

var ErrUnexpected = errors.New("unexpected error occured")
var ErrQuestionNotFound = errors.New("question not found")
var ErrQuestionAlreadyAnswered = errors.New("question already answered")
var ErrQuestionSessionNotRunning = errors.New("question session not running")
var ErrUserAlreadyVoted = errors.New("user already voted")
var ErrQuestionMaxLengthExceeded = errors.New("question max length exceeded")
var ErrQuestionNotOwned = errors.New("question not owned")
var ErrUserHasNotVoted = errors.New("user has not voted")
