package tokenprovider

import (
	"errors"
	"social-todo-list/common"
)

type Provider interface {
	Gernerate(data ToKenPayload, expiry int) (Token, error)
	Validate(token string) (ToKenPayload, error)
	SecretKey() string
}

type ToKenPayload interface {
	UserId() int
	Role() string
}

type Token interface {
	GetToken() string
}

var (
	ErrNotFound = common.NewCustomError(errors.New("token not found"),
		"token not found", "ErrNotFound")
	ErrEncodingToken = common.NewCustomError(errors.New("error encoding the token"),
		"error encoding the token", "ErrEncodingToken")
	ErrInvalidToken = common.NewCustomError(errors.New("invalid token provided"),
		"invalid token provided", "ErrInvalidToken")
)
