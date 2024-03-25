package usecases_events

import shared_models "voting/shared/models"

type UserBoundEvent struct {
	Event Event
	User  shared_models.UserContext
}
