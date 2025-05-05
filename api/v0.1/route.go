package api

import (
	"net/http"

	"github.com/Immortanbird/api_using_go/api/v0.1/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterMiddleware(g *gin.Engine, mw ...gin.HandlerFunc) {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(mw...)

	zap.L().Info("Middleware configured.")
}

func RegisterRouter(g *gin.Engine) {
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// Self check (starting with "/check/...")
	group := g.Group("/check")
	{
		group.GET("/ping", handler.Ping)
		group.GET("/health", handler.HealthCheck)
		group.GET("/cpu", handler.CPUCheck)
		group.GET("/ram", handler.RamCheck)
		group.GET("/disk", handler.DiskCheck)
	}

	zap.L().Info("Router configured.")
}
