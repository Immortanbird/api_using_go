package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg string, cfgPaths ...string) error {
	if cfg != "" {
		// Load the configuration file.
		i := strings.LastIndex(cfg, ".")
		viper.SetConfigName(cfg[:i])
		viper.SetConfigType(cfg[i+1:])
	} else {
		// Load the default configuration file.
		viper.SetConfigFile("./config/config.json")
	}

	for _, path := range cfgPaths {
		viper.AddConfigPath(path)
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			viper.SetConfigFile("config.yml")
			return viper.ReadInConfig()
		} else {
			// Config file was found but another error was produced
			return err
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	return nil
}

func OpenDB() (*gorm.DB, error) {
	conf := viper.GetStringMap("db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		conf["user"],
		conf["password"],
		conf["host"],
		conf["port"],
		conf["database"],
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Database connected.")

	return db, nil
}
