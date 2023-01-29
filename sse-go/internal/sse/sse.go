package sse

type EventType string

const (
	NEW_QUESTION    EventType = "new_question"
	UPVOTE_QUESTION EventType = "upvote_question"
	ANSWER_QUESTION EventType = "answer_question"
)

type Event struct {
	EventType
	Payload string
}

type NewQuestionEvent struct {
	Id   string
	Text string
}
