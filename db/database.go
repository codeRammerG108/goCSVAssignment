package db

import (
	models "gocsv/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=localhost user=postgres password=admin@123 dbname=testDatabase port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	// Auto-migrate the schema
	DB.AutoMigrate(&models.Employee{})
}
