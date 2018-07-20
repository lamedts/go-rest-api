package core

import (
	"go-rest-api/config"
	"os"
	"path"
	"path/filepath"
	"strings"
	// . "go-rest-api/core/syncer"
	// . "go-rest-api/core/tasks"

	// "go-rest-api/rpc/messaging"
	// . "go-rest-api/types"
	// . "go-rest-api/types/datastore"

	"go-rest-api/logger"
)

var debugMode = false
var coreLogger = logger.GetLogger("core")

func init() {
	zakkayaEnv := os.Getenv("APP_ENV")
	if strings.Contains(strings.ToUpper(zakkayaEnv), "DEBUG") {
		debugMode = true
	}
}

type Yay struct {
	config config.Config
}

var yayInstance *Yay

func GetYay(configuration *config.Config, versionString string) *Yay {
	coreLogger.Infof(".............    config    .............")
	// coreLogger.Infof("StoreCode    : %+v", configuration.StoreCode)
	// coreLogger.Infof("AppDataServer: %+v", configuration.AppDataServerSettings.Host)
	coreLogger.Infof("........................................")
	if yayInstance != nil {
		return yayInstance
	}

	/*
		Establish Postgresql database connection
		Initialize the database instance
		Run restart script (for checking)
	*/
	var zakkayaImplicitConfigPath string
	if debugMode {
		zakkayaImplicitConfigPath = path.Join("..", "datastore")
	} else {
		zakkayaImplicitConfigPath = path.Join(filepath.Dir(os.Args[0]), ".zakkaya")
	}
	if _, err := os.Stat(zakkayaImplicitConfigPath); os.IsNotExist(err) {
		if err := os.MkdirAll(zakkayaImplicitConfigPath, 0777); err != nil {
			coreLogger.Warnln("Failed to create directory", zakkayaImplicitConfigPath)
		}
	}
	if absPath, err := filepath.Abs(zakkayaImplicitConfigPath); err == nil {
		zakkayaImplicitConfigPath = absPath
	}

	// var dataStore DataStore
	// switch configuration.DatastoreSettings.Type {
	// // case DATASTORE_TYPE_SQLITE:
	// // 	coreLogger.Warnln("loading sqlite")
	// // 	workingDir, _ := os.Getwd()
	// // 	dataStore = datastore.NewSqliteDB(path.Join(workingDir, "/dummy.db"))
	// default:
	// 	pgConfig, err := config.LoadPgConfigFromJson(path.Join(zakkayaImplicitConfigPath, "db_config.json"))
	// 	if pgConfig == nil && err != nil {
	// 		coreLogger.Warn("Failed to load Postgres Config file:", err)
	// 		coreLogger.Warn("now use default Postgres Configuration")
	// 		pgConfig = &config.DefaultPGConfig
	// 		if err := pgConfig.SavePgConfigToJsonFile(path.Join(zakkayaImplicitConfigPath, "db_config.json")); err != nil {
	// 			coreLogger.Warnln("Failed to save pg config to json file", err)
	// 		}
	// 	}
	// 	pgConfig.ScriptFolder = path.Join(zakkayaImplicitConfigPath, "sql")
	// 	// restartScriptFile := path.Join(zakkayaImplicitConfigPath, "sql", "restart.sql")
	// 	// pgConfig.RestartScriptFile = restartScriptFile
	// 	dataStore = datastore.NewPgDB(*pgConfig)
	// }
	// if dataStore == nil {
	// 	coreLogger.Fatal("Failed to establish database connection")
	// 	return nil
	// }

	// if len(configuration.StoreCode) > 0 {
	// 	dataStore.SetStoreCode(configuration.StoreCode)
	// }

	yayInstance = &Yay{
		config: *configuration,
	}
	return yayInstance
}

func (yay *Yay) GetConfig() config.Config {
	return yay.config
}

func (yay *Yay) Start() {

	// for _, server := range yay.rpcServers {
	// 	if !server.Serve() {
	// 		coreLogger.Error("Failed to serve")
	// 	}
	// }
}
