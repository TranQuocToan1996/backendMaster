package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,alphanum"`
}

type renewResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != payload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expire session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		payload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
