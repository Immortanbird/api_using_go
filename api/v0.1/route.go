package api

import (
	"net/http"
	"time"

	"github.com/Immortanbird/api_using_go/api/v0.1/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterMiddleware(g *gin.Engine, mw ...gin.HandlerFunc) {
	// Define your CORS configuration
	config := cors.Config{
		// AllowOrigins: []string{"http://localhost:8080", "https://your-vue-app.com"}, // Frontend URLs
		AllowOrigins:     []string{"*"}, // Allows all origins. For production, list specific origins.
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"}, // Headers the browser is allowed to access
		AllowCredentials: true,                       // Important for cookies, authorization headers with HTTPS
		// AllowOriginFunc: func(origin string) bool { // For more complex origin checking
		//  return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour, // How long the results of a preflight request can be cached
	}

	g.Use(cors.New(config))

	zap.L().Info("CORS configured.")

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
