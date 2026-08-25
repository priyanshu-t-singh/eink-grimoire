package core

import (
	"fmt"
	"le-grimoire/internal/constants"
	"log/slog"
	"time"
)

type Config struct {
	Version   string
	Server    serverConfig
	Logs      logsConfig
	KavitaAPI kavitaAPIConfig
	Display   displayConfig
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
	Scheme     string
	APIKey     string
	PluginName string
	Timeout    time.Duration
}

type logsConfig struct {
	Dir string
}

// TODO: Implement a proper configuration loading mechanism (e.g., from a file or environment variables)
func NewConfig(logger *slog.Logger) *Config {
	logger.Info("Loading configuration...")

	defaultHost := constants.Host
	defaultPort := constants.Port

	return &Config{
		Version: constants.Version,
		Server: serverConfig{
			Host: defaultHost,
			Port: defaultPort,
		},
		Logs: logsConfig{
			Dir: constants.LogFileDirectory,
		},

		// Kavita Configuration
		KavitaAPI: kavitaAPIConfig{
			Scheme:     constants.KavitaScheme,
			Host:       constants.KavitaAPIHost,
			Port:       constants.KavitaAPIPort,
			PluginName: constants.KavitaAPIPluginName,
			Timeout:    constants.KavitaAPITimeout,
			APIKey:     constants.KavitaAPIKey,
		},

		// Display Configuration (4.2" e-paper = 400x300)
		Display: displayConfig{
			Width:  constants.DisplayWidth,
			Height: constants.DisplayHeight,
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
	return fmt.Sprintf("%s://%s:%d", cfg.KavitaAPI.Scheme, cfg.KavitaAPI.Host, cfg.KavitaAPI.Port)
}

func (cfg *Config) GetDisplayResolution() string {
	return fmt.Sprintf("%dx%d", cfg.Display.Width, cfg.Display.Height)
}
