package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	group := g.Group("/check")
	{
		group.GET("/ping", Ping)
		group.GET("/health", HealthCheck)
		group.GET("/cpu", CPUCheck)
		group.GET("/ram", RamCheck)
		group.GET("/disk", DiskCheck)
	}

	log.Println("Router configured.")
}
