package token

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestPasto_HAPPY(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal(err)
	}
	maker, err := NewPasetoMaker(config)
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expriedAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)
	require.True(t, payload.Username == username)
	require.WithinDuration(t, issuedAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expriedAt, payload.ExpiredAt, time.Second)
	require.True(t, payload.ExpiredAt.Sub(payload.IssueAt) == duration)
}

func TestPaseto_Expire(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal(err)
	}
	maker, err := NewPasetoMaker(config)
	require.NoError(t, err)

	username := util.RandomOwner()

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrTokenExpire.Error())
	require.Nil(t, payload)
}

func TestPaseto_Invalid(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := tokenObj.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(config.TokenSymetricKey))
	maker, err := NewPasetoMaker(config)
	require.NoError(t, err)

	getPayload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, getPayload)
}
