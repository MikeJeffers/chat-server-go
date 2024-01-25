package main

import "os"

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) < 1 {
		return defaultValue
	}
	return value
}

var (
	DB_USER     = getEnv("POSTGRES_USER", "user")
	DB_PASSWORD = getEnv("POSTGRES_PASSWORD", "password")
	DB_NAME     = getEnv("POSTGRES_DB", "test")
	DB_HOST     = getEnv("POSTGRES_HOST", "localhost")
	DB_PORT     = getEnv("POSTGRES_PORT", "5432")

	REDIS_HOST     = getEnv("REDIS_HOST", "localhost")
	REDIS_PORT     = getEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD = getEnv("REDIS_PASSWORD", "")

	SECRET_JWT = getEnv("SECRET_JWT", "idk")
)
