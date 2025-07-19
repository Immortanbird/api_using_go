package main

import (
	"net/http"
	"time"

	"github.com/Immortanbird/api_using_go/api/v0.1/handler"
	"github.com/Immortanbird/api_using_go/internal/config"
	"github.com/Immortanbird/api_using_go/internal/middlewares"
	"github.com/Immortanbird/api_using_go/internal/repository/crud"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config := config.LoadConfig()

	crud.OpenDB(config)

	g := gin.New()
	gin.SetMode((*config).App.Mode)
	setRouting(g, config)

	// Run the server
	g.Run((*config).App.Addr + ":" + (*config).App.Port)
}

func setRouting(g *gin.Engine, config *config.Config) {
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// Self check (starting with "/check/...")
	healthCheckRouters := g.Group("/check")
	{
		healthCheckRouters.GET("/ping", handler.Ping)
		healthCheckRouters.GET("/health", handler.HealthCheck)
		healthCheckRouters.GET("/cpu", handler.CPUCheck)
		healthCheckRouters.GET("/ram", handler.RamCheck)
		healthCheckRouters.GET("/disk", handler.DiskCheck)
	}

	handler := handler.Handler{Config: config}

	unprotected := g.Group("")
	unprotected.POST("/login", handler.Login)
	unprotected.POST("/register", handler.Register)
	unprotected.POST("/refresh-token", handler.Refresh)

	// Auth routes (starting with "/auth/...")
	userRouters := g.Group("")
	userRouters.DELETE("/user/delete", handler.DeleteAccount)

	imageRouters := g.Group("")
	imageRouters.POST("/image/upload", handler.UploadImage)
	imageRouters.GET("/image/download", handler.DownloadImage)
	imageRouters.DELETE("/image/delete", handler.DeleteImage)

	mw := middlewares.Middleware{JWTSecretKey: (*config).JWT.SecretKey}
	userRouters.Use(mw.Authenticate)
	imageRouters.Use(mw.Authenticate)

	// Define your CORS configuration
	g.Use(cors.New(cors.Config{
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
	}))
	g.Use(gin.Recovery())

	zap.L().Info("Routing configured.")
}
