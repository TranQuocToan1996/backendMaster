package api

import (
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func NewServer(store *db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
