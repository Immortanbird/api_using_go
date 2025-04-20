package repository

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("db.user"),
		os.Getenv("db.pwd"),
		os.Getenv("db.host"),
		os.Getenv("db.port"),
		os.Getenv("db.name"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Panic("Failed to connect to database.")
	}

	zap.L().Info("Database connected.")

	return db
}
