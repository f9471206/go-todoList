// services/main_test.go
package services_test

import (
	"log"
	"os"
	"testing"
	"todolist/config"
)

func TestMain(m *testing.M) {
	// 載入測試用環境變數檔案，並初始化config
	if err := config.LoadEnv(".env.test"); err != nil {
		log.Println("警告：載入 .env.test 失敗")
	}
	os.Setenv("GO_ENV", "test")

	code := m.Run()
	os.Exit(code)
}
