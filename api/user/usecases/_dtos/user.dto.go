package user_dtos

import "voting/shared/auth"

type UserDto struct {
	Email string
	Name  string
	Token string
	Role  auth.Role
}
