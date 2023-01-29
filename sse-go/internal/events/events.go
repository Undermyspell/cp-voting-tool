package events

type Events string

const (
	NEW_QUESTION    Events = "new_question"
	UPVOTE_QUESTION Events = "upvote_question"
	ANSWER_QUESTION Events = "answer_question"
)
