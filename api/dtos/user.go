package dtos

type NewTestUserDto struct {
	FirstName string `json:"text" binding:"required"`
	LastName  string `json:"anonymous"  binding:"required"`
}
