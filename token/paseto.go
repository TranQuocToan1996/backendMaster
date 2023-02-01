package token

import (
	"fmt"
	"time"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMakerV2 struct {
	paseto      paseto.Protocol
	symetricKey []byte
}

func NewPasetoMaker(config util.Config) (Maker, error) {
	if len(config.TokenSymetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size")
	}

	maker := &PasetoMakerV2{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(config.TokenSymetricKey),
	}

	return maker, nil
}

func (p *PasetoMakerV2) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return p.paseto.Encrypt(p.symetricKey, payload, nil)
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
