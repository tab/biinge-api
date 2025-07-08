package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const DebugLevel = "debug"

type AppConfig struct {
	Name        string
	Environment string
	LogLevel    string
	ClientURL   string
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type TMDBConfig struct {
	BaseURL            string
	BaseImageURL       string
	APIReadAccessToken string
	Locale             string
	Timeout            time.Duration
}

type Config struct {
	App         AppConfig
	Server      ServerConfig
	DatabaseDSN string
	JWT         JWTConfig

	TMDBConfig
}

func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", env),
		fmt.Sprintf(".env.%s.local", env),
	}
	for _, file := range envFiles {
		_ = godotenv.Overload(file)
	}

	return &Config{
		App: AppConfig{
			Environment: env,
			Name:        getEnvString("APP_NAME"),
			ClientURL:   getEnvString("CLIENT_URL"),
			LogLevel:    getEnvString("LOG_LEVEL"),
		},
		Server: ServerConfig{
			Address:      getEnvString("SERVER_ADDRESS"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},

		DatabaseDSN: getEnvString("DATABASE_DSN"),

		JWT: JWTConfig{
			SecretKey:            getEnvString("JWT_SECRET_KEY"),
			AccessTokenDuration:  24 * time.Hour,
			RefreshTokenDuration: 7 * 24 * time.Hour,
		},

		TMDBConfig: TMDBConfig{
			BaseURL:            getEnvString("TMDB_BASE_URL"),
			BaseImageURL:       getEnvString("TMDB_BASE_IMAGE_URL"),
			APIReadAccessToken: getEnvString("TMDB_API_READ_ACCESS_TOKEN"),
			Locale:             getEnvString("TMDB_LOCALE"),
		},
	}
}

func getEnvString(envVar string) string {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}

	return ""
}
