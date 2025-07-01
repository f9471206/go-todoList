package main

import (
	"log"
	"todolist/config"
	"todolist/db/seed"
)

func main() {
	// 載入環境變數
	if err := config.LoadEnv(".env.local"); err != nil {
		if err := config.LoadEnv(".env"); err != nil {
			log.Fatal("❌ 無法載入任何環境變數檔案")
		}
	}

	config.ConnectDatabase()

	seed.Seed(config.DB)

	log.Println("✅ 種子資料執行完成")
}
