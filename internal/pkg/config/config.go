package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	Port                 int
	GitlabToken          string
	SlackToken           string
	SlackFallbackUser    string
	FilteredStartings    []string
	FilteredEmails       []string
	FilteredGroups       []string
	CommitLogType        string
	CommitLogServer      string
	CommitLogServicename string
	DatabasePath         string
}

// Init does the configuration init
func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		return errors.Wrap(err, "Unable to read config file")
	}

	return nil
}

// Load loads the configuration from viper and returns a Config instance
func Load() *Config {
	port := viper.GetInt("Server.Port")
	gitlabToken := viper.GetString("Gitlab.Token")
	slackToken := viper.GetString("Slack.Token")
	slackFallbackUser := viper.GetString("Slack.FallbackUser")
	startings := viper.GetStringSlice("Filters.Startings")
	emails := viper.GetStringSlice("Filters.Emails")
	groups := viper.GetStringSlice("Filters.Groups")
	commitLogType := viper.GetString("CommitLog.Type")
	commitLogServer := viper.GetString("CommitLog.Server")
	commitLogSericename := viper.GetString("CommitLog.Servicename")
	dbPath := viper.GetString("Database.Path")

	return &Config{
		port,
		gitlabToken,
		slackToken,
		slackFallbackUser,
		startings,
		emails,
		groups,
		commitLogType,
		commitLogServer,
		commitLogSericename,
		dbPath,
	}
}
