package shared

import shared_models "voting/shared/models"

type EventType string

type Event struct {
	EventType
	Payload string
}

type UserBoundEvent struct {
	Event Event
	User  shared_models.UserContext
}
