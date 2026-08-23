package server

import (
	"le-grimoire/internal/core"
)

func StartServer() {
	startApp()
}

func startApp() {
	// Create the app instance
	app := core.NewApp()
	app.InitLogging()

	// Initialize the HTTP server
	svr := core.NewHTTPServer(app)

	// Initialize the routes
	// TODO: Implement route initialization logic here

	// Run the server
	core.RunHTTPServer(app, svr)

	// TODO: Implement graceful shutdown logic here
}
