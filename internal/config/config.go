package config

import (
	"log"
	"os"
)

// Config holds all application configuration.
// Values are read from environment variables.
type Config struct {
	MongoURI    string
	MongoDBName string
	JWTSecret   string
	ServerPort  string
	GinMode     string
}

// Load reads configuration from environment variables.
// Call this once at startup.
func Load() *Config {
	cfg := &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://devops:w3eL5SnEx245I48f7McN@10.107.20.93:27017/IMZ_USER?directConnection=true&authSource=admin"),
		MongoDBName: getEnv("MONGO_DB_NAME", "gst-api"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-min-32-chars"),
		ServerPort:  getEnv("SERVER_PORT", "8520"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("[CONFIG] JWT_SECRET environment variable is required but not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
