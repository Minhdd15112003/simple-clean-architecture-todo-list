package jwt

import (
	"fmt"
	"social-todo-list/common"
	"social-todo-list/components/tokenprovider"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secret string
	prefix string
}

func NewTokenJWTProvider(prefix string) *jwtProvider {
	return &jwtProvider{
		prefix: prefix,
	}
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}

func (t *jwtProvider) SecretKey() string {
	return t.secret
}

func (j *jwtProvider) Gernerate(data tokenprovider.ToKenPayload, expiry int) (tokenprovider.Token, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodES256, myClaims{
		common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiry) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})
	myToken, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return nil, err
	}

	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (tokenprovider.ToKenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}
	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}
	claims, ok := res.Claims.(*myClaims)

	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}
	return claims.Payload, nil
}
