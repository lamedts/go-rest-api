package logger

import (
	"go-rest-api/config"
	"go-rest-api/core/global"
	"io"
	"os"
	"path"
	"strings"

	"io/ioutil"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logLevel logrus.Level
var primaryOutStream io.Writer
var fileFormatter = prefixed.TextFormatter{FullTimestamp: true, ForceFormatting: true}

var (
	appLogHook   logrus.Hook
	warnLogHook  logrus.Hook
	errorLogHook logrus.Hook

	coreLogHook       logrus.Hook
	datastoreLogHook  logrus.Hook
	serverLogHook     logrus.Hook
	apiRequestLogHook logrus.Hook
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
	serverLogHook = lfshook.NewHook(
		getPathMap(path.Join(config.LogPath, "server.log")),
		&fileFormatter,
	)
	apiRequestLogHook = lfshook.NewHook(
		getPathMap(path.Join(config.LogPath, "apiRequest.log")),
		&fileFormatter,
	)

}

func setPrimaryOutStream() {
	logLevel = logrus.InfoLevel
	primaryOutStream = ioutil.Discard // abandone output
	if global.DebugMode {
		logLevel = logrus.DebugLevel
		primaryOutStream = os.Stdout
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
	case "server", "api":
		logger.AddHook(serverLogHook)
	case "apiRequest":
		logger.AddHook(apiRequestLogHook)
	}

	return logger.WithField("prefix", module)
}
