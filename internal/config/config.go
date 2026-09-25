package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"tarmo/internal/lib/logger"

	"github.com/joho/godotenv"
)

var Env = "development"
var Version = "dev"

type Config struct {
	Port           string
	DBPath         string
	AllowedOrigins []string
}

func Load() *Config {
	_ = godotenv.Load()

	var dbPath string
	if Env == "production" {
		dbPath = defaultDBPath()
	} else {
		dbPath = getEnv("DB_PATH", defaultDBPath())
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		logger.Fatal("could not create data directory: %v", err)
	}

	return &Config{
		Port:           getEnv("PORT", "9136"),
		DBPath:         dbPath,
		AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}
}

func defaultDBPath() string {
	if Env == "production" {
		return filepath.Join(appDataDir(), "data", "db.sqlite")
	}
	return "data/db.sqlite"
}

func appDataDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("ProgramData"), "Tarmo")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "share", "tarmo")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
