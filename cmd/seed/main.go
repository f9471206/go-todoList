package main

import (
	"log"
	"todolist/config"
	"todolist/db/seed"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("❌ 無法載入 .env 檔案")
	}

	config.ConnectDatabase()

	seed.Seed(config.DB)

	log.Println("✅ 種子資料執行完成")
}
