package server

import (
	"le-grimoire/internal/core"
	"le-grimoire/internal/handlers"
)

func StartServer() {
	startApp()
}

func startApp() {
	// Create the app instance
	app := core.NewApp()
	app.InitLogging()

	// Initialize the HTTP server
	router := core.NewHTTPServer()

	// Initialize the routes
	handlers.InitRoutes(app, router)

	// Run the server
	core.RunHTTPServer(app, router)

	// TODO: Implement graceful shutdown logic here
}
