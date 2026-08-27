package constants

import (
	"time"
)

const (
	Host = "127.0.0.1"
	Port = "8080"
)

// Injected at build time using ldflags, default is 0.1.0 for local development
var Version = "0.1.0"

const (
	KavitaAPIBaseURL    = "http://kavita:5000"
	KavitaAPIPluginName = "le-grimoire"
	KavitaAPITimeout    = 10 * time.Second
)

// Display Configuration (4.2" e-paper = 400x300)
const (
	DisplayHeight = 300
	DisplayWidth  = 400
)

const (
	SqliteDatabasePath = "./.db/database.sqlite"
	LogFileDirectory   = "./.logs"
)
