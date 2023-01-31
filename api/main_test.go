package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
