package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDB(config *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DbSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	log.Println("Successfully connected to database")
	return db
}
