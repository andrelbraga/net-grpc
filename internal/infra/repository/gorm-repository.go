package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=dbmaster dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Print("Database connected!")
	return db, nil
}
