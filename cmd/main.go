package main

import (
	"github.com/Immortanbird/api_using_go/config"
	"github.com/Immortanbird/api_using_go/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.Init("config.json", "json", "./config")

	// var logger *zap.Logger
	// switch viper.GetString("mode") {
	// case "debug":
	// 	logger, _ = zap.NewDevelopment()
	// case "test":
	// 	logger = zap.NewExample()
	// case "prod":
	// 	logger, _ = zap.NewProduction()
	// default:
	// 	logger, _ = zap.NewProduction()
	// }
	logger, _ := zap.NewDevelopment()
	// defer logger.Sync()

	logger.Info("logger created", zap.String("sugar", "no"))

	// g := gin.Default()
	g := gin.New()
	gin.SetMode(viper.GetString("mode"))

	handler.Load(g)

	g.Run("localhost:" + viper.GetString("port"))
}
