package logger

import (
	"go-rest-api/config"
	"io"
	"os"
	"path"
	"strings"

	"io/ioutil"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logLevel logrus.Level
var primaryOutStream io.Writer
var fileFormatter = logrus.TextFormatter{FullTimestamp: true}

var (
	appLogHook   logrus.Hook
	warnLogHook  logrus.Hook
	errorLogHook logrus.Hook

	coreLogHook      logrus.Hook
	datastoreLogHook logrus.Hook
)

func init() {
	setPrimaryOutStream()

	// General hooks
	appLogHook = lfshook.NewHook(
		getPathMap(path.Join(config.LogPath, path.Base(os.Args[0])+".log")),
		&fileFormatter,
	)

	warnLogHook = lfshook.NewHook(
		lfshook.PathMap{
			logrus.WarnLevel:  path.Join(config.LogPath, "warn.log"),
			logrus.ErrorLevel: path.Join(config.LogPath, "warn.log"),
			logrus.FatalLevel: path.Join(config.LogPath, "warn.log"),
			logrus.PanicLevel: path.Join(config.LogPath, "warn.log"),
		},
		&fileFormatter,
	)

	errorLogHook = lfshook.NewHook(
		lfshook.PathMap{
			logrus.ErrorLevel: path.Join(config.LogPath, "error.log"),
			logrus.FatalLevel: path.Join(config.LogPath, "error.log"),
			logrus.PanicLevel: path.Join(config.LogPath, "error.log"),
		},
		&fileFormatter,
	)

	// Module hooks
	coreLogHook = lfshook.NewHook(
		getPathMap(path.Join(config.LogPath, "core.log")),
		&fileFormatter,
	)
	datastoreLogHook = lfshook.NewHook(
		getPathMap(path.Join(config.LogPath, "datastore.log")),
		&fileFormatter,
	)

}

func setPrimaryOutStream() {
	appEnv := os.Getenv("APP_ENV")
	switch strings.ToUpper(appEnv) {
	case "DEBUG", "TEST":
		logLevel = logrus.DebugLevel
		primaryOutStream = os.Stdout
	default:
		logLevel = logrus.InfoLevel
		primaryOutStream = ioutil.Discard // abandone output
	}
}

func getPathMap(logPath string) lfshook.PathMap {
	return lfshook.PathMap{
		logrus.DebugLevel: logPath,
		logrus.InfoLevel:  logPath,
		logrus.WarnLevel:  logPath,
		logrus.ErrorLevel: logPath,
		logrus.FatalLevel: logPath,
		logrus.PanicLevel: logPath,
	}
}

// GetLogger to let other package get their own logger
func GetLogger(module string) *logrus.Entry {
	logger := logrus.New()
	logger.Formatter = &fileFormatter
	logger.SetLevel(logLevel)
	setPrimaryOutStream()

	logger.AddHook(appLogHook)
	logger.AddHook(warnLogHook)
	logger.AddHook(errorLogHook)

	switch strings.ToLower(module) {
	case "core", "main":
		logger.AddHook(coreLogHook)
	case "datastore":
		logger.AddHook(datastoreLogHook)
	}

	return logger.WithField("prefix", module)
}
