package services

func NewTestUser() UserService {
	return &TestUserService{}
}
