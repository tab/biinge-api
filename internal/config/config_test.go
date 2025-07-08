package config

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"biinge-api/pkg/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *Config
	}{
		{
			name: "Success",
			args: []string{},
			env:  map[string]string{},
			expected: &Config{
				App: AppConfig{
					Environment: "test",
					Name:        "biinge",
					ClientURL:   "http://localhost:3000",
					LogLevel:    "info",
				},
				Server: ServerConfig{
					Address: "localhost:8080",
				},
				DatabaseDSN:   "postgres://postgres:postgres@localhost:5432/biinge-test?sslmode=disable",
				SecretKeyBase: "SECRET",
				JWTSecretKey:  "SECRET",
				TMDBConfig: TMDBConfig{
					BaseURL:            "https://api.themoviedb.org/3",
					BaseImageURL:       "https://image.tmdb.org/t/p",
					APIReadAccessToken: "SECRET",
					Locale:             "en-US",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			result := LoadConfig()

			assert.Equal(t, tt.expected.App.Environment, result.App.Environment)
			assert.Equal(t, tt.expected.App.Name, result.App.Name)
			assert.Equal(t, tt.expected.App.ClientURL, result.App.ClientURL)
			assert.Equal(t, tt.expected.App.LogLevel, result.App.LogLevel)
			assert.Equal(t, tt.expected.Server.Address, result.Server.Address)
			assert.Equal(t, tt.expected.DatabaseDSN, result.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKeyBase, result.SecretKeyBase)
			assert.Equal(t, tt.expected.JWTSecretKey, result.JWTSecretKey)
			assert.Equal(t, tt.expected.TMDBConfig.BaseURL, result.TMDBConfig.BaseURL)
			assert.Equal(t, tt.expected.TMDBConfig.BaseImageURL, result.TMDBConfig.BaseImageURL)
			assert.Equal(t, tt.expected.TMDBConfig.APIReadAccessToken, result.TMDBConfig.APIReadAccessToken)
			assert.Equal(t, tt.expected.TMDBConfig.Locale, result.TMDBConfig.Locale)
			assert.Equal(t, tt.expected.TMDBConfig.Timeout, result.TMDBConfig.Timeout)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
