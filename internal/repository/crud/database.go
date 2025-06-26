package crud

import (
	"fmt"

	"github.com/Immortanbird/api_using_go/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB(config *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		(*config).DB.User,
		(*config).DB.Pwd,
		(*config).DB.Host,
		(*config).DB.Port,
		(*config).DB.Name,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Panic("Failed to connect to database: " + err.Error())
	}

	zap.L().Info("Database connected.")
}
