package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"todolist/config"
	"todolist/db/seed"
	"todolist/middleware"
	"todolist/routes"
	"todolist/utils"

	_ "todolist/docs"
)

// ✅ 提早載入 .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ 無法載入 .env 檔案")
	}
}

func main() {
	utils.InitLogger(true)
	utils.Logger.Info("Logger 初始化成功")

	config.ConnectDatabase()

	// ✅ 呼叫種子資料
	seed.Seed(config.DB)

	r := gin.Default()
	r.Use(middleware.RecoveryMiddleware())
	routes.RegisterRoutes(r)

	if config.ENV != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Run(":" + config.AppPort)
}
