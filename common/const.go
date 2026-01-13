package common

import (
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	hasher, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hasher), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

const (
	CurrentUser = "current_user"
)
