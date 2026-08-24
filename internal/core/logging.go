package core

import (
	"fmt"
	"le-grimoire/internal/util"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
)

func (a *App) InitLogging() {
	consoleHandler := tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	if a.Config.Logs.Dir != "" {
		if err := os.MkdirAll(a.Config.Logs.Dir, 0755); err != nil {
			log.Fatalf("failed to create logs directory: %v\n", err)
		}
	}

	logFilePath := filepath.Join(a.Config.Logs.Dir, fmt.Sprintf("le-grimoire-%s.log", time.Now().Format("2006-01-02_15-04-05")))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to create log file: %v\n", err)
	}

	jsonHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	// Combine into root logger
	multiHandler := util.NewMultiHandler(consoleHandler, jsonHandler)
	wrappedHandler := &util.ContextHandler{Handler: multiHandler}
	logger := slog.New(wrappedHandler)

	// Bind to App & Global Default
	a.Logger = logger
	slog.SetDefault(logger)

	// Set up shutdown function for logger
	a.ShutdownLogger = func() {
		slog.Info("Shutting down logger...")
		if err := logFile.Close(); err != nil {
			slog.Error("failed to close log file", "error", err)
		}
	}
}
