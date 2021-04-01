package ability

import (
	"github.com/iamolegga/enviper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ReadConfiguration() *Configuration {
	conf, err := readConfiguration()
	if err != nil { // Handle errors reading the config file
		logrus.WithError(err).Fatalf("Error reading config.")
	} else {
		logrus.Infof("The configuration has been successfully ridden.")
		logrus.Debugf("-> %+v", conf)
	}
	return conf
}

func readConfiguration() (*Configuration, error) {
	e := enviper.New(viper.New())
	e.SetEnvPrefix("ABILITY")

	e.SetConfigName("config")
	e.SetConfigType("toml")
	e.AddConfigPath("/etc/ability/")
	e.AddConfigPath("$HOME/.ability")
	e.AddConfigPath(".")
	if err := e.ReadInConfig(); err != nil {
		return nil, err
	}

	var C Configuration
	//_ = e.Unmarshal(&C)
	if err := e.Unmarshal(&C); err != nil {
		return nil, err
	}
	return &C, nil
}
