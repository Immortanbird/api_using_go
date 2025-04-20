package main

import (
	"os"

	"github.com/Immortanbird/api_using_go/config"
	"github.com/Immortanbird/api_using_go/handler"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	config.LoadEnv()
	config.LoadLogger()
}

func main() {
	// g := gin.Default()
	g := gin.New()
	gin.SetMode(os.Getenv("mode"))

	handler.Load(g)

	g.Run(os.Getenv("addr") + ":" + os.Getenv("port"))
}
