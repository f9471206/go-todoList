package config

import "os"

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

func init() {
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
