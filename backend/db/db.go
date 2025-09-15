package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Connect(models ...interface{}) {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is empty")
	}
	var err error
	Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	if err := Conn.AutoMigrate(models...); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}
}
