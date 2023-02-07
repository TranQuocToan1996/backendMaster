package token

import (
	"testing"
	"time"

	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expriedAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)
	require.True(t, payload.Username == username)
	require.WithinDuration(t, issuedAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expriedAt, payload.ExpiredAt, time.Second)
}

func TestJWTMaker_Expire(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()

	token, _, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrTokenExpire.Error())
	require.Nil(t, payload)
}

func TestJWTMaker_Invalid(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	token, err := tokenObj.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	getPayload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, getPayload)
}
