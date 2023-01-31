package api

import (
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func NewServer(store db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
