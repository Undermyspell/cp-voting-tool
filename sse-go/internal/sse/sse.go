package sse

type EventType string

const (
	NEW_QUESTION    EventType = "new_question"
	UPVOTE_QUESTION EventType = "upvote_question"
	ANSWER_QUESTION EventType = "answer_question"
	RESET_SESSION   EventType = "reset_session"
)

type Event struct {
	EventType
	Payload string
}

const PayloadEmpty = "{}"
