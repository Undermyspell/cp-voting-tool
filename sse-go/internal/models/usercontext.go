package models

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type UserContext struct {
	Name  string
	Email string
}

func (userContext *UserContext) GetHash() string {
	marshalled, _ := json.Marshal(userContext)
	bytes, _ := bcrypt.GenerateFromPassword(marshalled, 10)
	return string(bytes)
}

const User string = "user"
