package models

import (
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sse/internal/models/roles"

	"golang.org/x/crypto/sha3"
)

type UserContext struct {
	Name  string
	Email string
	Role  roles.Role
}

func (userContext *UserContext) GetHash(secret string) string {
	h := hmac.New(sha3.New256, []byte(secret))
	marshalled, _ := json.Marshal(userContext)
	h.Write(marshalled)
	hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprint(hash)
}

const User string = "user"
