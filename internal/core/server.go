package core

import (
	"le-grimoire/internal/middleware"
	"net/http"
)

func NewHTTPServer() *http.ServeMux {
	return http.NewServeMux()
}

func RunHTTPServer(app *App, router *http.ServeMux) {
	server := &http.Server{
		Addr: app.Config.GetServerAddr(),
		Handler: middleware.CreateStack(
			middleware.Logging,
		)(router),
	}

	serverAddr := app.Config.GetServerAddr()
	app.Logger.Info("Starting HTTP server", "address", serverAddr)
	app.Logger.Info("Server is running at " + app.Config.GetServerURI())
	app.Logger.Info("Display resolution: " + app.Config.GetDisplayResolution())

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error("HTTP server error", "error", err)
	}
}
