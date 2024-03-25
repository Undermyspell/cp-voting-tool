package usecases_events

import "voting/shared/shared_models"

type UserBoundEvent struct {
	Event Event
	User  shared_models.UserContext
}
