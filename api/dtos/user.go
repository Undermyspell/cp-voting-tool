package dtos

type NewTestUserDto struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname"  binding:"required"`
}
