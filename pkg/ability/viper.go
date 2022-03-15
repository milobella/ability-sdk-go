package ability

import (
	"fmt"
	"strings"

	"github.com/iamolegga/enviper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const configName = "ability"

// ReadConfiguration
// Deprecated: use ReadConfig instead, which allows to specify a richer configuration.
func ReadConfiguration() *Configuration {
	var config Configuration
	ReadConfig(&config)
	return &config
}

func ReadConfig(config *Configuration) {
	e := enviper.New(viper.New())
	e.SetEnvPrefix(strings.ToUpper(configName))

	e.SetConfigName("config")
	e.SetConfigType("toml")
	e.AddConfigPath(fmt.Sprintf("/etc/%s/", configName))
	e.AddConfigPath(fmt.Sprintf("$HOME/.%s", configName))
	e.AddConfigPath(".")
	err := e.ReadInConfig()
	if err != nil {
		fatal(err)
	}

	if err = e.Unmarshal(&config); err != nil {
		fatal(err)
	} else {
		logrus.Info("Successfully red configuration !")
		logrus.Debugf("-> %+v", config)
	}

	logrus.SetLevel(config.Server.LogLevel)
}

func fatal(err error) {
	logrus.WithError(err).Fatal("Error reading config.")
}
