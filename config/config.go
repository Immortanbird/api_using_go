package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadLogger() {
	var logger *zap.Logger

	switch os.Getenv("mode") {
	case "debug":
		logger, _ = zap.NewDevelopment()
	case "test":
		logger = zap.NewExample()
	case "release":
		logger, _ = zap.NewProduction()
	default:
		panic("Mode unknown. Available mode: debug release test")
	}

	defer logger.Sync()

	// Replace the global logger, so that it can be used elsewhere
	defer zap.ReplaceGlobals(logger)

	zap.L().Info("logger created", zap.String("sugar", "no"))
}
