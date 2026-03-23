package config

import (
	"os"
	"strconv"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	// Port is the HTTP port the server listens on.
	Port string

	// RedisAddr is the Redis server address (host:port).
	RedisAddr string

	// RedisPassword is the optional Redis authentication password.
	RedisPassword string

	// RedisDB is the Redis database index to use.
	RedisDB int

	// CacheTTLMinutes is the number of minutes RSS articles are cached in Redis.
	CacheTTLMinutes int

	// MaxArticlesPerFeed is the maximum number of articles fetched per RSS feed.
	MaxArticlesPerFeed int

	// FetchTimeoutSeconds is the HTTP timeout in seconds for RSS feed requests.
	FetchTimeoutSeconds int
}

// Load reads configuration from environment variables and returns a Config.
// Each variable falls back to a sensible default if not set.
func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "3001"),
		RedisAddr:           getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", ""),
		RedisDB:             getEnvInt("REDIS_DB", 0),
		CacheTTLMinutes:     getEnvInt("CACHE_TTL_MINUTES", 15),
		MaxArticlesPerFeed:  getEnvInt("MAX_ARTICLES_PER_FEED", 5),
		FetchTimeoutSeconds: getEnvInt("FETCH_TIMEOUT_SECONDS", 10),
	}
}

// getEnv returns the value of an environment variable,
// or the provided fallback if the variable is not set.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getEnvInt returns the integer value of an environment variable,
// or the provided fallback if the variable is not set or cannot be parsed.
func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
