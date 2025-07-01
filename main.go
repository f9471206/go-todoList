package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"todolist/config"
	"todolist/middleware"
	"todolist/routes"
	"todolist/utils"

	_ "todolist/docs"
)

func main() {
	// 明確載入環境變數檔案
	// 嘗試載入 .env.local，失敗再載 .env
	if err := config.LoadEnv(".env.local"); err != nil {
		if err := config.LoadEnv(".env"); err != nil {
			log.Fatal("❌ 無法載入任何環境變數檔案")
		}
	}

	utils.InitLogger(true)
	utils.Logger.Info("Logger 初始化成功")

	config.ConnectDatabase()

	r := gin.Default()
	r.Use(middleware.RecoveryMiddleware())
	routes.RegisterRoutes(r)

	if config.ENV != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	for _, route := range r.Routes() {
		fmt.Printf("Route: %s %s\n", route.Method, route.Path)
	}

	r.Run(":" + config.AppPort)
}
