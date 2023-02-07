package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TranQuocToan1996/backendMaster/model"
	"github.com/TranQuocToan1996/backendMaster/token"
	"github.com/gin-gonic/gin"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(model.AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("authorization header is not corrected")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if !strings.EqualFold(model.AuthorizationTypeBearer, fields[0]) {
			err := fmt.Errorf("authorization type is not Bearer, got: %v", fields[0])
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(model.AuthorizationPayloadKey, payload)
		c.Next()
	}
}
