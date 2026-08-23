package core

import (
	"fmt"
	"le-grimoire/internal/constants"
	"log/slog"
)

type Config struct {
	Version string
	Server  struct {
		Host string
		Port int
	}
	Logs struct {
		Dir string
	}
	KavitaAPI struct {
		Host string
		Port int
	}
}

type serverConfig struct {
	Host string
	Port int
}

type logsConfig struct {
	Dir string
}

// TODO: Implement a proper configuration loading mechanism (e.g., from a file or environment variables)
func NewConfig(logger *slog.Logger) *Config {
	logger.Debug("Loading configuration...")

	defaultHost := "127.0.0.1"
	defaultPort := 8080

	return &Config{
		Version: constants.Version,
		Server: serverConfig{
			Host: defaultHost,
			Port: defaultPort,
		},
		Logs: logsConfig{
			Dir: "./.logs",
		},
		KavitaAPI: serverConfig{
			Host: defaultHost,
			Port: 5000,
		},
	}
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}

func (cfg *Config) GetServerURI() string {
	return fmt.Sprintf("http://%s", cfg.GetServerAddr())
}
