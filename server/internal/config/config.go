package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	OpenAIKey   string
	PostgresURL string
	RedisURL    string
	ServerPort  string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		OpenAIKey:   mustGet("OPENAI_API_KEY"),
		PostgresURL: mustGet("POSTGRES_URL"),
		RedisURL:    mustGet("REDIS_URL"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env: %s", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
