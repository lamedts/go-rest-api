package types

const (
	DefaultConfigFilename = "config.yaml"
	DefaultAPIPort        = 8080
)

type DatastoreConfig struct {
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
type GoogleMapAPICOnfig struct {
	Key string `yaml:"key"`
}

type ServerConfig struct {
	APIServerConfig     *APIServerConfig     `yaml:"api"`
	GraphqlServerConfig *GraphqlServerConfig `yaml:"graphql"`
}
type utilityConfig struct {
	GoogleMapAPIConfig *GoogleMapAPICOnfig `yaml:"googlemapapi"`
}

type Config struct {
	DatastoreSettings DatastoreConfig `yaml:"db"`
	ServerSettings    ServerConfig    `yaml:"server"`
	UtilitySettings   utilityConfig   `yaml:"utility"`
}
