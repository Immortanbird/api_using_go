package repository

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
