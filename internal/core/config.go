package core

import (
	"fmt"
	"le-grimoire/internal/constants"
	"le-grimoire/internal/util"
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
	BaseURL    string
	APIKey     string
	PluginName string
	Timeout    time.Duration
}

type logsConfig struct {
	Dir string
}

func NewConfig(logger *slog.Logger) *Config {
	logger.Info("Loading configuration...")

	return &Config{
		Version: constants.Version,
		Server: serverConfig{
			Host: util.GetServerHost(),
			Port: util.GetServerPort(),
		},
		Logs: logsConfig{
			Dir: util.GetLogFileDirectory(),
		},

		// Kavita Configuration
		KavitaAPI: kavitaAPIConfig{
			BaseURL:    util.GetKavitaBaseURL(),
			PluginName: constants.KavitaAPIPluginName,
			Timeout:    constants.KavitaAPITimeout,
			APIKey:     util.GetKavitaAPIKey(),
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
	return cfg.KavitaAPI.BaseURL
}

func (cfg *Config) GetDisplayResolution() string {
	return fmt.Sprintf("%dx%d", cfg.Display.Width, cfg.Display.Height)
}
