package core

import (
	"database/sql"
	"le-grimoire/internal/constants"
	"le-grimoire/internal/kavita"
	"le-grimoire/internal/util"
	"log/slog"
)

type App struct {
	// Config
	Config   *Config
	Logger   *slog.Logger
	Database *sql.DB
	Version  string

	// Kavita Repository
	KavitaRepository *kavita.Repository

	// Shutdown
	ShutdownLogger func()
}

func NewApp() *App {
	logger := util.NewLogger()
	logger.Info("Initializing application...")

	return &App{
		Config:           NewConfig(logger),
		Logger:           logger,
		Version:          constants.Version,
		KavitaRepository: &kavita.Repository{},
	}
}
