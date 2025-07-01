package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JWTSecret string
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	AppPort   string
	ENV       string
)

// LoadEnv 載入指定的 env 檔案，並設定全局變數
func LoadEnv(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		log.Printf("❌ 無法載入 %s 檔案: %v", filename, err)
		return err
	}
	loadEnvVars()
	return nil
}

// loadEnvVars 從 os.Getenv 讀取設定並賦值給全局變數
func loadEnvVars() {
	JWTSecret = mustGetenv("JWT_SECRET")
	DBUser = mustGetenv("DB_USER")
	DBPass = mustGetenv("DB_PASSWORD")
	DBHost = mustGetenv("DB_HOST")
	DBPort = mustGetenv("DB_PORT")
	DBName = mustGetenv("DB_NAME")
	AppPort = mustGetenv("APP_PORT")
	ENV = mustGetenv("GO_ENV")
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Environment variable " + key + " not set")
	}
	return val
}

func IsLocal() bool {
	return ENV != "production"
}
