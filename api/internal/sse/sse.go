package sse

type EventType string

const (
	NEW_QUESTION    EventType = "new_question"
	UPVOTE_QUESTION EventType = "upvote_question"
	ANSWER_QUESTION EventType = "answer_question"
	STOP_SESSION    EventType = "stop_session"
	START_SESSION   EventType = "start_session"
)

type Event struct {
	EventType
	Payload string
}

const PayloadEmpty = "{}"
