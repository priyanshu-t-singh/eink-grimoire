package constants

import (
	"time"
)

const (
	DefaultHost = "127.0.0.1"
	DefaultPort = 8321
)

// NOTE: Injected at build time using ldflags, default is `dev` for local development
var Version = "dev"

const (
	KavitaAPIBaseURL    = "http://kavita:5000"
	KavitaAPIPluginName = "le-grimoire"
	KavitaAPITimeout    = 10 * time.Second
)

// TODO: Make this configurable for each device type, and/or make it auto-detectable
// Display Configuration (4.2" e-paper = 400x300)
const (
	DisplayHeight = 300
	DisplayWidth  = 400
)

// TODO: move it to .config/le-grimoire/config.toml and make it configurable
// TODO: get the config path from env or default to $HOME/.config/le-grimoire/config.toml
const (
	ConfigFileName = "config.toml"
)
