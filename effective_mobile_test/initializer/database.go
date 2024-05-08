package initializer

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maghavefun/effective_mobile_test/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Migrate() {
	err := DB.AutoMigrate(&model.Person{}, &model.Car{})
	if err != nil {
		log.Fatal("Failed to run migrations", err)
	}
}

func ConnectToDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	f, err := os.Create("gorm.log")
	if err != nil {
		log.Fatal("Error creating log file for gorm:", err)
	}

	DB.Logger = logger.New(
		log.New(f, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)
}
