package core

import (
	"le-grimoire/internal/constants"
	"le-grimoire/internal/kavita"
	"le-grimoire/internal/util"
	"log/slog"
)

type App struct {
	// Config
	Config  *Config
	Logger  *slog.Logger
	Version string

	// Kavita Repository
	KavitaRepository *kavita.Repository
}

func NewApp() *App {
	logger := util.NewLogger()

	return &App{
		Config:           NewConfig(logger),
		Logger:           logger,
		Version:          constants.Version,
		KavitaRepository: &kavita.Repository{},
	}
}
