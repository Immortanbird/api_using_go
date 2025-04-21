package main

import (
	"github.com/Immortanbird/api_using_go/config"
	"github.com/Immortanbird/api_using_go/models"
	"go.uber.org/zap"
)

func main() {
	config.LoadEnv()
	config.LoadLogger()

	models.OpenDB()
	models.Migrate()

	zap.L().Info("Database migration completed.")
}
