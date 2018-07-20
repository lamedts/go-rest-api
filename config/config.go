package config

import (
	"go-rest-api/types/datastore"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var LogPath = path.Join(".", "log") //path.Join("/", "var", "log", "yay")
var yayLoggerFile *os.File

func init() {
	logPath := path.Join(".", "log") //path.Join("/", "var", "log", "yay")
	appEnv := os.Getenv("APP_ENV")
	if runtime.GOOS != "linux" || strings.Contains(strings.ToUpper(appEnv), "DEBUG") {
		logPath = path.Join(".", "log")
	} else {
		logPath = path.Join("/", "var", "log", "yay")
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			log.Warnln(err)
		} else {
			log.Infoln("Directory", logPath, "is now created")
		}
	}

	switch strings.ToUpper(appEnv) {
	case "DEBUG", "TEST":
	default:
		fullLogPath := path.Join(logPath, path.Base(os.Args[0])+".log")
		f, err := os.OpenFile(fullLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Warnln("Error opening file", err)
		} else {
			// log.Infoln("Writing logs to", fullLogPath)
		}
		yayLoggerFile = f
	}
}

func GetYayFile() *os.File {
	return yayLoggerFile
}

const (
	DefaultConfigFilename string = "config.yaml"
)

type DatastoreSettings struct {
	Type string `yaml:"type"`
}

type apiServerSettings struct {
	Port int `yaml:"port"`
}

type Config struct {
	DatastoreSettings DatastoreSettings `yaml:"db"`
	APIServerSettings apiServerSettings `yaml:"api_server"`
}

var DefaultConfig Config = Config{
	DatastoreSettings: DatastoreSettings{
		Type: datastore.DatastoreTypeMysql,
	},
	APIServerSettings: apiServerSettings{
		Port: 8080,
	},
}

func LoadConfigFromYaml(filename string) *Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Cannot read config file '", filename, "'")
		return nil
	}
	var config = &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Error("Cannot parse config file '", filename, "'")
		return nil
	}
	return config
}
func (config *Config) SaveConfigToYamlFile(filename string) error {
	if yamlBytes, err := yaml.Marshal(config); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, yamlBytes, 0644)
	}
}
