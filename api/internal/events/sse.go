package events

import "voting/internal/models"

type UserBoundSseEvent struct {
	Event Event
	User  models.UserContext
}
