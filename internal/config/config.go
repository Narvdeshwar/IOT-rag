package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey   string
	GeminiKey   string
	PostgresURL string
	RedisURL    string
	ServerPort  string
	OllamaURL   string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		GeminiKey:   getEnv("GEMINI_API_KEY", ""),
		PostgresURL: mustGet("POSTGRES_URL"),
		RedisURL:    mustGet("REDIS_URL"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434"),
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
