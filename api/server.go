package api

import (
	"fmt"

	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/TranQuocToan1996/backendMaster/token"
	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func (s *Server) Start() error {
	return s.router.Run(s.config.ServerAddress)
}

func (s *Server) setupRouter() {
	router := gin.Default()
	s.router = router

	router.POST("/users", s.createUser)
	router.POST("/users/login", s.loginUser)

	authGroup := router.Group("/").Use(authMiddleware(s.tokenMaker))
	authGroup.POST("/accounts", s.createAccount)
	authGroup.GET("/accounts/:id", s.getAccount)
	authGroup.GET("/accounts", s.listAccount)
	authGroup.POST("/transfers", s.createTransfer)
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config)
	if err != nil {
		return nil, fmt.Errorf("cant create token marker %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
