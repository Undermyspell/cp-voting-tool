package sse

import "sse/internal/models"

type EventType string

const (
	NEW_QUESTION    EventType = "new_question"
	UPDATE_QUESTION EventType = "update_question"
	UPVOTE_QUESTION EventType = "upvote_question"
	ANSWER_QUESTION EventType = "answer_question"
	DELETE_QUESTION EventType = "delete_question"
	STOP_SESSION    EventType = "stop_session"
	START_SESSION   EventType = "start_session"
)

type Event struct {
	EventType
	Payload string
}

type UserBoundSseEvent struct {
	Event Event
	User  models.UserContext
}

const PayloadEmpty = "{}"
