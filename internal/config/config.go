package config

import (
	"strings"
	"tickets/internal/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port     int
		GRPCPort int
	}
	Database struct {
		URL      string
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Logging logger.Config `json:"logging" yaml:"logging"`
	Mode    string
	Port    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Enable environment variable support
	viper.AutomaticEnv()

	// Set environment variable key mappings
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Map environment variables to config keys
	if err := viper.BindEnv("server.port", "SERVER_PORT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("server.grpc_port", "SERVER_GRPC_PORT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.host", "DATABASE_HOST"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.port", "DATABASE_PORT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.user", "DATABASE_USER"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.password", "DATABASE_PASSWORD"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.dbname", "DATABASE_DBNAME"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.url", "DATABASE_URL"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("logging.level", "LOGGING_LEVEL"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("logging.format", "LOGGING_FORMAT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("logging.output", "LOGGING_OUTPUT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("logging.include_caller", "LOGGING_INCLUDE_CALLER"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("logging.include_timestamp", "LOGGING_INCLUDE_TIMESTAMP"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("mode", "MODE"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("port", "PORT"); err != nil {
		return nil, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set defaults if not provided
	if cfg.Mode == "" {
		cfg.Mode = "debug"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Server.GRPCPort == 0 {
		cfg.Server.GRPCPort = 9090
	}

	return &cfg, nil
}
