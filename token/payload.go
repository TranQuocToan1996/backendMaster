package token

import (
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenExpire  = errors.New("token expire")
	ErrInvalidToken = errors.New("token invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssueAt   time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *Payload) Valid() error {
	if time.Now().UnixNano() > p.ExpiredAt.UnixNano() {
		return ErrTokenExpire
	}

	if reflect.DeepEqual(&Payload{}, p) {
		return ErrInvalidToken
	}

	return nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil

}
