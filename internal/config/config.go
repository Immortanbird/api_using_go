package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Check  CheckConfig  `mapstructure:"check"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	PicBed PicBedConfig `mapstructure:"pic_bed"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Addr string `mapstructure:"addr"`
	Port string `mapstructure:"port"`
}

type CheckConfig struct {
	MaxPingCount int `mapstructure:"max_ping_count"`
}

type DBConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"pwd"`
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SecretKey  string `mapstructure:"secret_key"`
	ExpRefresh int    `mapstructure:"refresh_lifespan"`
	ExpAccess  int    `mapstructure:"access_lifespan"`
}

type PicBedConfig struct {
	GitHubToken     string `mapstructure:"github_token"`
	GitHubUrlFormat string `mapstructure:"github_url_format"`
}

func LoadConfig() *Config {
	// Load environment variables
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config: %w", err))
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	// Create logger
	var logger *zap.Logger

	switch viper.GetString("app.mode") {
	case "debug":
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(fmt.Errorf("fatal error creating logger: %w", err))
		}
	case "test":
		logger = zap.NewExample()
	case "release":
		logger, err = zap.NewProduction()
		if err != nil {
			panic(fmt.Errorf("fatal error creating logger: %w", err))
		}
	default:
		panic("Mode unknown. Available mode: debug release test")
	}

	// Replace the global logger, so that it can be used elsewhere
	zap.ReplaceGlobals(logger)

	zap.L().Info("Logger created:", zap.String("sugar", "no"))

	return config
}
