package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	Server struct {
		LogLevel string
		Port     int
	}
	Gitlab struct {
		Token string
	}
	CommitLog struct {
		Type        string
		Server      string
		Servicename string
	}
	Slack struct {
		Emoji        string
		Token        string
		FallbackUser string
	}
	Filters struct {
		Startings []string
		Email     []string
		Groups    []string
	}
	Database struct {
		Path string
	}
}

// Load loads the configuration from viper and returns a Config instance
func Load(cfgFile string) (*Config, error) {
	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig()

	if err != nil {
		return nil, errors.Wrap(err, "Unable to read config file")
	}
	var conf Config

	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to unmarshal config")
	}

	return &conf, nil
}