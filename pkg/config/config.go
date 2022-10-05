package config

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Server Server          `mapstructure:"server"`
	Tools  map[string]Tool `mapstructure:"tools"`
}

// fun String() : Serialization function of Config (for logging)
func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Config serialization error %s", err)
	}
	return string(b)
}

type Server struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type Tool struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
