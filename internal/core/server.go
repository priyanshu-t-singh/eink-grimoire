package core

import (
	"le-grimoire/internal/middleware"
	"net/http"
)

func NewHTTPServer(app *App) *http.Server {
	router := http.NewServeMux()

	server := &http.Server{
		Addr: app.Config.GetServerAddr(),
		Handler: middleware.CreateStack(
			middleware.Logging,
		)(router),
	}

	return server
}

func RunHTTPServer(app *App, server *http.Server) {
	serverAddr := app.Config.GetServerAddr()
	app.Logger.Info("Starting HTTP server", "address", serverAddr)
	app.Logger.Info("Server is running at " + app.Config.GetServerURI())

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error("HTTP server error", "error", err)
	}
}
