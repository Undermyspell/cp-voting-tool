package events

import "voting/internal/models"

type UserBoundEvent struct {
	Event Event
	User  models.UserContext
}
