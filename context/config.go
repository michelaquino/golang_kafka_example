package context

import (
	"os"
	"strconv"
	"sync"
)

// LogConfig represents the log configuration
type LogConfig struct {
	LogLevel    string
	LogToFile   bool
	LogFileName string
}

// APIConfig represents the API configuration
type APIConfig struct {
	LogConfig *LogConfig
}

var apiConfig *APIConfig
var onceConfig sync.Once

// GetAPIConfig return the instance of the APIConfig
func GetAPIConfig() *APIConfig {
	onceConfig.Do(func() {
		apiConfig = &APIConfig{
			LogConfig: getLogConfig(),
		}
	})

	return apiConfig
}

func getLogConfig() *LogConfig {
	return &LogConfig{
		LogLevel:    os.Getenv("LOG_LEVEL"),
		LogToFile:   getLogToFileFlag(),
		LogFileName: os.Getenv("LOG_FILE_NAME"),
	}
}

func getLogToFileFlag() bool {
	if logToFile, err := strconv.ParseBool(os.Getenv("LOG_TO_FILE")); err == nil {
		return logToFile
	}

	return false
}
