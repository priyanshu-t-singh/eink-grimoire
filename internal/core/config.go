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
		Host       string
		Port       int
		PluginName string
	}
	Display struct {
		Height int
		Width  int
	}
}

type serverConfig struct {
	Host string
	Port int
}

type displayConfig struct {
	Height int
	Width  int
}

type kavitaAPIConfig struct {
	Host       string
	Port       int
	PluginName string
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

		// Kavita Configuration
		KavitaAPI: kavitaAPIConfig{
			Host:       defaultHost,
			Port:       5000,
			PluginName: "le-grimoire",
		},

		// Display Configuration (4.2" e-paper = 400x300)
		Display: displayConfig{
			Width:  400,
			Height: 300,
		},
	}
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}

func (cfg *Config) GetServerURI() string {
	return fmt.Sprintf("http://%s", cfg.GetServerAddr())
}

func (cfg *Config) GetKavitaAPIURI() string {
	return fmt.Sprintf("http://%s:%d", cfg.KavitaAPI.Host, cfg.KavitaAPI.Port)
}

func (cfg *Config) GetDisplayResolution() string {
	return fmt.Sprintf("%dx%d", cfg.Display.Width, cfg.Display.Height)
}
