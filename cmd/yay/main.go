package main

import (
	"os"
	"path"
	"time"

	"fmt"

	"go-rest-api/config"
	. "go-rest-api/core"
	"go-rest-api/logger"

	"github.com/coreos/go-systemd/daemon"
	"github.com/sirupsen/logrus"
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
var mainLogger *logrus.Entry = logger.GetLogger("main")

func init() {
	if len(BUILD_DATE) == 0 {
		BUILD_DATE = "Now --> " + time.Now().Local().Format(time.RFC3339)
	}
	cobra.OnInitialize()
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "zakkaya config yaml file")
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
		os.Exit(zerror.EXITCODE_PROGRAM_COMMAND_ERROR)
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
	mainLogger.Infoln("Starting zakkaya version:", versionString)
	configuration := loadConfig()
	zakkaya := zakkayaSetup(configuration, versionString)

	daemon.SdNotify(false, "READY=1")

	mainLogger.Infoln("Finished zakkaya Configuration and Setup")
	mainLogger.Infoln("----------------------------------------")
	zakkaya.Start()
	mainLogger.Infoln("-------------------------")
	mainLogger.Infoln("zakkaya Start Serving Now")
	mainLogger.Infoln("-------------------------")

	<-quit
}

var quit = make(chan bool)

func loadConfig() *config.Config {
	var configuration config.Config
	useDefaultConfig := true
	if len(cfgFile) > 0 {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			useDefaultConfig = true
		} else {
			mainLogger.Infoln("Loading zakkaya configuration from", cfgFile)
			if tmpConfig := config.LoadConfigFromYaml(cfgFile); tmpConfig == nil {
				useDefaultConfig = true
			} else {
				configuration = *tmpConfig
				useDefaultConfig = false
			}
		}
	} else {
		workingDir, _ := os.Getwd()
		defaultConfigFile := path.Join(workingDir, config.DefaultConfigFilename)
		if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
			useDefaultConfig = true
		} else {
			if tmpConfig := config.LoadConfigFromYaml(defaultConfigFile); tmpConfig == nil {
				useDefaultConfig = true
			} else {
				configuration = *tmpConfig
				useDefaultConfig = false
			}
		}
	}

	if useDefaultConfig {
		mainLogger.Infoln("Loading default zakkaya configuration")
		configuration = config.DEFAULT_ZAKKAYA_CONFIG
		if workingDir, err := os.Getwd(); err != nil {
			mainLogger.Errorln("Failed to identify current working directory", err)
			os.Exit(yayerror.EXITCODE_UNEXPECTED_ERROR)
		} else {
			if err := configuration.SaveConfigToYamlFile(path.Join(workingDir, config.DefaultConfigFilename)); err != nil {
				mainLogger.Warnln("Failed to save config file to " + path.Join(workingDir, config.DefaultConfigFilename))
			}
		}
	}
	return &configuration
}

func yaySetup(configuration *config.Config, versionString string) *Yay {
	if configuration != nil {

		if _, err := yaml.Marshal(*configuration); err == nil {
			// mainLogger.Infoln("Successful load configuration")
			// mainLogger.Debugf("Configuration:")
			// mainLogger.Debugf(string(bytes))
		}

		// zakkaya := GetZakkaya(configuration, versionString)
		return GetYay(configuration, versionString)
	}
	return nil
}
