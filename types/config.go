package types

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	DefaultConfigFilename = "config.yaml"
	DefaultAPIPort        = 8080
)

type DatastoreSettings struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	DBname   string `yaml:"dbname"`
}

type APIServerConfig struct {
	Port int `yaml:"port"`
}
type GraphqlServerConfig struct {
}

type ServerConfig struct {
	APIServerConfig     *APIServerConfig     `yaml:"api"`
	GraphqlServerConfig *GraphqlServerConfig `yaml:"graphql"`
}

type Config struct {
	DatastoreSettings DatastoreSettings `yaml:"db"`
	Servers           ServerConfig      `yaml:"server"`
	// APIServerSettings apiServerSettings `yaml:"api_server"`
}

func (config *Config) SaveConfigToYamlFile(filename string) error {
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, yamlBytes, 0644)
}
