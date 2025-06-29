package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 依據環境決定資料庫名稱（可擴充）
	dbName := DBName
	if ENV != "production" {
		// 例如開發或測試環境，連另一個資料庫（你可改成你想要的）
		dbName = DBName
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPass, DBHost, DBPort, dbName)

	var err error

	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to database:", dbName)
			return
		}
		log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to database after multiple attempts:", err)
}
