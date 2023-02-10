package models

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sse/internal/models/roles"
)

type UserContext struct {
	Name  string
	Email string
	Role  roles.Role
}

func (userContext *UserContext) GetHash() string {
	marshalled, _ := json.Marshal(userContext)

	algorithm := fnv.New32a()
	algorithm.Write(marshalled)
	hash := algorithm.Sum32()
	return fmt.Sprint(hash)
}

const User string = "user"
