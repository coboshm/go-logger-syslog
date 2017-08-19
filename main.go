package main

import (
	"flag"
	"fmt"

	"github.com/coboshm/go-logger-syslog/logger"
	"github.com/spf13/viper"
)

// Handle CLI arguments.
var env = flag.String("env", "testing", "Set application environment: -env=testing")
var configDir = flag.String("config-dir", "config", "Set the application config dir: -config=/opt/appname/current/config")

func init() {
	flag.Parse()
}

func main() {
	configFileName := "config"
	appName := "appTestName"

	v := viper.GetViper()

	v.AddConfigPath(*configDir)
	v.SetConfigName(configFileName)

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to load config from %s", *configDir))
	}

	loggerDSN := viper.GetString("log.dsn")
	log, err := logger.NewLoggerFromDSN(loggerDSN, appName, *env)
	if err != nil {
		panic(fmt.Sprintf("Error creating new logger: %s", err))
	}

	log.Info("Running...", logger.NewField("newField1", "value1"), logger.NewField("newField2", 2))
	log.Debug("Debugging Running...", logger.NewField("newField1", "value1"), logger.NewField("newField2", 2))
}
