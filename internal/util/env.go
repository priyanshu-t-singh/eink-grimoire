package util

import (
	"le-grimoire/internal/constants"
	"os"
)

func Getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
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
