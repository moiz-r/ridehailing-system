package configs

import (
	"log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	GRPCPort    string         `mapstructure:"grpcPort"`
	Database    DatabaseConfig `mapstructure:"database"`
	MetricsPort string         `mapstructure:"metricsPort"`
}

// LoadConfig reads configuration from a file or environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/") // Path to look for the config file
	viper.AutomaticEnv()

	// Read the configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	// Map configuration to struct
	var config Config
	if err := viper.UnmarshalKey("config", &config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
		return nil, err
	}

	return &config, nil
}
