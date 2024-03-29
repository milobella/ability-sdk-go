package config

import (
	"fmt"
	"strings"

	"github.com/iamolegga/enviper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const configName = "ability"

func Read(extensions ...interface{}) *Config {
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

	var config Config
	if err = e.Unmarshal(&config); err != nil {
		fatal(err)
	} else {
		logrus.Info("Successfully red configuration !")
		logrus.Debugf("-> %+v", config)
	}

	for _, ext := range extensions {
		if err = e.Unmarshal(ext); err != nil {
			fatal(err)
		}
	}

	var level logrus.Level
	level, err = logrus.ParseLevel(config.Server.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.WithError(err).Errorf("Error occured while parsing log level. Default to INFO")
	} else {
		logrus.SetLevel(level)
	}

	return &config
}

func fatal(err error) {
	logrus.WithError(err).Fatal("Error reading config.")
}
