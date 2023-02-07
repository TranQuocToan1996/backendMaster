package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/TranQuocToan1996/backendMaster/model"
	"github.com/TranQuocToan1996/backendMaster/token"
	"google.golang.org/grpc/metadata"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeaders := md.Get(model.AuthorizationHeaderKey)
	if len(authHeaders) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := authHeaders[0]
	fields := strings.Fields(authHeader)

	isValidBearer := func(fields []string) bool {
		return len(fields) >= 2
	}

	if !isValidBearer(fields) {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := fields[0]
	if strings.EqualFold(authType, model.AuthorizationTypeBearer) {
		return nil, fmt.Errorf("unsupport authorization type")
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token %s", err)
	}

	return payload, nil
}
