package core

import (
	"database/sql"
	"le-grimoire/internal/constants"
	"le-grimoire/internal/device"
	"le-grimoire/internal/kavita"
	"le-grimoire/internal/render"
	"le-grimoire/internal/state"
	"le-grimoire/internal/util"
	"log/slog"
)

// TODO: set up a proper dependency injection framework for this app,
// so that we can easily swap out implementations of the repositories, renderers, etc.
type App struct {
	// Config
	Config   *Config
	Logger   *slog.Logger
	Database *sql.DB
	Version  string

	// Repository
	KavitaRepository kavita.Repository
	DeviceRepository *device.Repository

	StateMachine *state.Machine
	Renderer     *render.Renderer
	FrameCache   *render.FrameCache

	// Shutdown
	ShutdownLogger func()
}

func NewApp() *App {
	logger := util.NewLogger()
	logger.Info("Initializing application...")

	return &App{
		Config:   NewConfig(logger),
		Logger:   logger,
		Database: nil,
		Version:  constants.Version,
	}
}
