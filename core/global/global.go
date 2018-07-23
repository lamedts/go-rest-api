package global

import (
	"os"
	"strings"
)

var DebugMode = false

func init() {
	appEnv := os.Getenv("APP_ENV")
	if strings.Contains(strings.ToUpper(appEnv), "DEBUG") {

		DebugMode = true
	}
}
