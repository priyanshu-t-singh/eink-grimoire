package util

import (
	"le-grimoire/internal/constants"
	"log/slog"
	"os"
	"strconv"
)

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetServerHost() string {
	return Getenv("SERVER_HOST", constants.Host)
}

func GetServerPort() int {
	port, err := strconv.Atoi(Getenv("SERVER_PORT", constants.Port))
	if err != nil {
		slog.Warn("Invalid SERVER_PORT value, using default port 8080", "warn", err)
		return 8080
	}
	return port
}

func GetKavitaBaseURL() string {
	return Getenv("KAVITA_BASE_URL", constants.KavitaAPIBaseURL)
}

func GetKavitaAPIKey() string {
	return Getenv("KAVITA_API_KEY", "")
}

func GetChromeRemoteURL() string {
	return Getenv("CHROME_REMOTE_URL", "")
}

func GetChromePath() string {
	return Getenv("CHROME_PATH", "")
}

func GetSqliteDatabasePath() string {
	return Getenv("SQLITE_DATABASE_PATH", constants.SqliteDatabasePath)
}

func GetLogFileDirectory() string {
	return Getenv("LOG_FILE_DIRECTORY", constants.LogFileDirectory)
}
