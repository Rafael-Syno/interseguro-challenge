// Package config centraliza la lectura de variables de entorno.
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	JWTSecret     string
	JWTTTL        time.Duration
	ClientID      string
	ClientSecret  string
	StatsAPIURL   string
	StatsTimeout  time.Duration
	RoundDecimals int
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "interChallSecret2000"),
		JWTTTL:        time.Duration(getEnvInt("JWT_TTL_MINUTES", 60)) * time.Minute,
		ClientID:      getEnv("CLIENT_ID", "interseguro"),
		ClientSecret:  getEnv("CLIENT_SECRET", "challenge2001"),
		StatsAPIURL:   getEnv("STATS_API_URL", "http://localhost:3000/api/v1/statistics"),
		StatsTimeout:  time.Duration(getEnvInt("STATS_TIMEOUT_SECONDS", 10)) * time.Second,
		RoundDecimals: getEnvInt("ROUND_DECIMALS", 10),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
