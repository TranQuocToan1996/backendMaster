package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	minSecretKeysize = 32
)

type JWTMaker struct {
	secretKey string
}

func (j *JWTMaker) CreateToken(username string, 
	duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	//TODO: Change alg RSA, elliptic digital signature
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := tokenObj.SignedString([]byte(j.secretKey))

	return token, payload, err
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var (
		signKeyProvider = func(t *jwt.Token) (interface{}, error) {
			//TODO: Change alg, elliptic digital signature
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidToken
			}
			return []byte(j.secretKey), nil
		}
		payload = &Payload{}
	)

	_, err := jwt.ParseWithClaims(token, payload, signKeyProvider)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrTokenExpire) {
			return nil, ErrTokenExpire
		}
		return nil, ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeysize {
		return nil, fmt.Errorf("key size at least %d chars", minSecretKeysize)
	}
	return &JWTMaker{secretKey}, nil
}
