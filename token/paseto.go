package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMakerV2 struct {
	paseto      paseto.Protocol
	symetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size")
	}

	maker := &PasetoMakerV2{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (p *PasetoMakerV2) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := p.paseto.Encrypt(p.symetricKey, payload, nil)
	return token, payload, err
}

func (p *PasetoMakerV2) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
