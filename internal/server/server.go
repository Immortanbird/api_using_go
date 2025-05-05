package server

import (
	"os"

	api "github.com/Immortanbird/api_using_go/api/v0.1"
	"github.com/Immortanbird/api_using_go/internal/config"
	"github.com/Immortanbird/api_using_go/internal/database"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	// Load environment
	config.LoadEnv()

	// Create logger
	config.LoadLogger()

	// Connect database
	database.OpenDB()

	// Create gin engine
	g := gin.New()
	gin.SetMode(os.Getenv("mode"))

	// Config routers
	api.RegisterMiddleware(g)
	api.RegisterRouter(g)

	// Run the server
	g.Run(os.Getenv("addr") + ":" + os.Getenv("port"))
}
