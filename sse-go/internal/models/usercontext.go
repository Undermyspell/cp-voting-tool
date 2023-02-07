package models

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
)

type UserContext struct {
	Name  string
	Email string
}

func (userContext *UserContext) GetHash() string {
	marshalled, _ := json.Marshal(userContext)

	algorithm := fnv.New32a()
	algorithm.Write(marshalled)
	hash := algorithm.Sum32()
	return fmt.Sprint(hash)
}

const User string = "user"
