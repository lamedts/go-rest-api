package main

import (
	"go-rest-api/core/global"
	"go-rest-api/types"
	"os"
	"path"
	"time"

	"fmt"

	"go-rest-api/config"
	"go-rest-api/core"
	yayerror "go-rest-api/errors"
	"go-rest-api/logger"

	"github.com/coreos/go-systemd/daemon"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

const softwareName = "yay"

var (
	BUILD_TAGS         = ""
	BUILD_DATE         string
	VERSION_MAJOR      = "0"
	VERSION_MINOR      = "0"
	VERSION_RELEASE    = "0"
	VERSION_DERIVATIVE = "unknown"
)

var cfgFile string
var mainLogger = logger.GetLogger("main")

func init() {
	if len(BUILD_DATE) == 0 {
		BUILD_DATE = "Now --> " + time.Now().Local().Format(time.RFC3339)
	}
	cobra.OnInitialize()
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config yaml file")
}

var RootCmd = &cobra.Command{
	Use:   softwareName,
	Short: softwareName + " core",
	Long:  softwareName + " core",
	Run: func(cmd *cobra.Command, args []string) {
		yay()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + softwareName,
	Run: func(cmd *cobra.Command, args []string) {
		if len(VERSION_DERIVATIVE) != 0 {
			fmt.Println(softwareName)
			fmt.Println("Version   :", VERSION_MAJOR+"."+VERSION_MINOR+"."+VERSION_RELEASE+"_"+VERSION_DERIVATIVE)
			fmt.Println("Build     :", BUILD_TAGS)
			fmt.Println("Build Date:", BUILD_DATE)
		} else {
			fmt.Println(softwareName)
			fmt.Println("Version   :", VERSION_MAJOR+"."+VERSION_MINOR+"."+VERSION_RELEASE)
			fmt.Println("Build     :", BUILD_TAGS)
			fmt.Println("Build Date:", BUILD_DATE)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print chosen configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := loadConfig()

		if _, err := yaml.Marshal(*configuration); err != nil {
			fmt.Println("Failed to parse the configuration")
		} else {
			// fmt.Println(softwareName + " configuration:")
			// fmt.Println(string(bytes))
			// fmt.Println("Successful load configuration")
		}

	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(yayerror.EXITCODE_PROGRAM_COMMAND_ERROR)
	}
}

func yay() {
	versionString := VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_RELEASE
	if len(VERSION_DERIVATIVE) != 0 {
		versionString = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_RELEASE + "_" + VERSION_DERIVATIVE
	}
	mainLogger.Infoln("----------------------------------------")
	mainLogger.Infoln("****************************************")
	mainLogger.Infoln("----------------------------------------")
	if global.DebugMode {
		mainLogger.Warnf("Starting app in debug mode")
	} else {
		mainLogger.Infoln("Starting version:", versionString)
	}
	configuration := loadConfig()
	yay := yaySetup(configuration)

	daemon.SdNotify(false, "READY=1")

	mainLogger.Infoln("Finished Configuration and Setup")
	mainLogger.Infoln("----------------------------------------")
	yay.Start()
	mainLogger.Infoln("-------------------------")
	mainLogger.Infoln("Start Serving Now")
	mainLogger.Infoln("-------------------------")

	<-quit
}

var quit = make(chan bool)

func loadConfig() *types.Config {
	var configuration types.Config
	workingDir, _ := os.Getwd()
	defaultConfigFile := path.Join(workingDir, types.DefaultConfigFilename)
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		mainLogger.Errorf("No Configuartion file in: %+v", defaultConfigFile)
		os.Exit(yayerror.EXITCODE_CONFIG_FILE_ERROR)
	} else {
		if tmpConfig := config.LoadConfigFromYaml(defaultConfigFile); tmpConfig == nil {
			mainLogger.Errorf("Cant load configuartion file")
			os.Exit(yayerror.EXITCODE_CONFIG_FILE_ERROR)
		} else {
			configuration = *tmpConfig
		}
	}
	return &configuration
}

func yaySetup(configuration *types.Config) *core.Yay {
	if configuration != nil {

		if _, err := yaml.Marshal(*configuration); err == nil {
			// mainLogger.Infoln("Successful load configuration")
			// mainLogger.Debugf("Configuration:")
			// mainLogger.Debugf(string(bytes))
		}

		// yay := GetYay(configuration, versionString)
		return core.GetYay(configuration)
	}
	return nil
}
