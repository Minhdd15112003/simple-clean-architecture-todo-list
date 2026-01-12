package common

import (
	"crypto/md5"
	"encoding/hex"
)

type TokenPayload struct {
	UId   int    `json:"user_id"`
	URole string `json:"role"`
}

func (payload TokenPayload) UserId() int {
	return payload.UId
}

func (payload TokenPayload) Role() string {
	return payload.URole
}

type md5Hash struct{}

func NewMd5Hash() *md5Hash {
	return &md5Hash{}
}

func (h *md5Hash) Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
