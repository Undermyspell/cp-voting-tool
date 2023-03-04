package models

import (
	"crypto/hmac"
	"encoding/hex"
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

func (userContext *UserContext) GetHash(secret string) string {
	h := hmac.New(fnv.New128a, []byte(secret))
	marshalled, _ := json.Marshal(userContext)
	h.Write(marshalled)
	hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprint(hash)
}

const User string = "user"
