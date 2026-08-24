package core

import (
	"context"
	"le-grimoire/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewHTTPServer() *http.ServeMux {
	return http.NewServeMux()
}

func RunHTTPServer(app *App, router *http.ServeMux) {
	server := &http.Server{
		Addr: app.Config.GetServerAddr(),
		Handler: middleware.CreateStack(
			middleware.Logging,
			middleware.AllowCors,
		)(router),
	}

	runServerWithGracefulShutdown(app, server)
}

func runServerWithGracefulShutdown(app *App, server *http.Server) {
	serverAddr := app.Config.GetServerAddr()
	app.Logger.Info("Starting HTTP server", "address", serverAddr)

	go func() {
		app.Logger.Info("Server is running at " + app.Config.GetServerURI())
		app.Logger.Info("Display resolution: " + app.Config.GetDisplayResolution())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("server failed to serve", "error", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Error("server forced to shutdown", "error", err)
	}

	app.ShutdownDatabase()
	app.ShutdownLogger()
	slog.Info("Shutdown complete")
}
