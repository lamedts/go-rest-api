package config

import (
	"go-rest-api/core/global"
	"go-rest-api/types"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var LogPath = path.Join(".", "log")
var yayLoggerFile *os.File

func init() {
	if runtime.GOOS != "linux" || global.DebugMode {
		LogPath = path.Join(".", "log")
	} else {
		LogPath = path.Join("/", "var", "log", "yay")
	}
	if _, err := os.Stat(LogPath); os.IsNotExist(err) {
		if err := os.MkdirAll(LogPath, 0755); err != nil {
			log.Warnln(err)
		} else {
			log.Infoln("Directory", LogPath, "is now created")
		}
	}

	// switch strings.ToUpper(appEnv) {
	// case "DEBUG", "TEST":
	// default:
	// 	fullLogPath := path.Join(LogPath, path.Base(os.Args[0])+".log")
	// 	f, err := os.OpenFile(fullLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	// 	if err != nil {
	// 		log.Warnln("Error opening file", err)
	// 	} else {
	// 		// log.Infoln("Writing logs to", fullLogPath)
	// 	}
	// 	yayLoggerFile = f
	// }
}

func GetYayFile() *os.File {
	return yayLoggerFile
}

func LoadConfigFromYaml(filename string) *types.Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Cannot read config file '", filename, "'")
		return nil
	}
	var config = &types.Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Error("Cannot parse config file '", filename, "'")
		return nil
	}
	return config
}
