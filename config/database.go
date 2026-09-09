package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"funfillers/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	var dsn string

	if databaseURL != "" {
		dsn = databaseURL
	} else {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "postgres"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "funfillers_db"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			host, user, password, dbname, port, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("⚠️ PostgreSQL connection failed: %v", err)
		log.Println("ℹ️ Backend will use fallback mock storage.")
		return nil
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("✅ Successfully connected to PostgreSQL database!")

	// Auto migrate database schemas
	if err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.SubCategory{},
		&models.Order{},
		&models.OrderItem{},
		&models.User{},
		&models.UIBanner{},
		&models.WishlistItem{},
		&models.Address{},
		&models.HeroVideoConfig{},
	); err != nil {
		log.Printf("⚠️ AutoMigrate error: %v\n", err)
	}

	DB = db
	return db
}
