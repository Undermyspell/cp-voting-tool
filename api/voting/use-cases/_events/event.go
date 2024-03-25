package usecases_events

import "voting/shared"

const (
	NEW_QUESTION         shared.EventType = "new_question"
	UPDATE_QUESTION      shared.EventType = "update_question"
	UPVOTE_QUESTION      shared.EventType = "upvote_question"
	UNDO_UPVOTE_QUESTION shared.EventType = "undo_upvote_question"
	ANSWER_QUESTION      shared.EventType = "answer_question"
	DELETE_QUESTION      shared.EventType = "delete_question"
	STOP_SESSION         shared.EventType = "stop_session"
	START_SESSION        shared.EventType = "start_session"
	USER_CONNECTED       shared.EventType = "user_connected"
	USER_DISCONNECTED    shared.EventType = "user_disconnected"
	HEART_BEAT           shared.EventType = "heart_beat"
)
