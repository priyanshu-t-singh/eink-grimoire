package constants

import "time"

var Host = "127.0.0.1" // Injected at build time using ldflags, default is localhost for local development

const (
	Version = "0.1.0"
	Port    = 8080
)

const (
	KavitaAPIPluginName = "le-grimoire"
	KavitaScheme        = "http"
	KavitaAPIHost       = "127.0.0.1"
	KavitaAPIKey        = "your_kavita_api_key_here" // Replace with your actual Kavita API key
	KavitaAPIPort       = 5000
	KavitaAPITimeout    = 10 * time.Second
)

// Display Configuration (4.2" e-paper = 400x300)
const (
	DisplayHeight = 300
	DisplayWidth  = 400
)

const SqliteDatabasePath string = "./.db/database.sqlite"
const LogFileDirectory string = "./.logs"
