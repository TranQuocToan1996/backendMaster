package gapi

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/TranQuocToan1996/backendMaster/pb"
	"github.com/TranQuocToan1996/backendMaster/token"
	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to get password %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exist %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to create user %s", err)
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "can not found user %s", err)
		}
		return nil, status.Errorf(codes.Internal, "fail to find user %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "incorrect password %s", err)
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create access token %s", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create refresh token %s", err)
	}

	mtdt := s.extractMetadata(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create session %s", err)
	}

	resp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return resp, nil
}

func (s *Server) mustEmbedUnimplementedSimpleBankServer() {}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cant create token marker %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
