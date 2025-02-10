package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config *Config
	mu     sync.RWMutex
)

type Config struct {
	LogLevel string `mapstructure:"log-level"`
	Server   struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Kafka struct {
		Brokers  []string `mapstructure:"brokers"`
		Topic    string   `mapstructure:"topic"`
		GroupID  string   `mapstructure:"group-id"`
		Offset   string   `mapstructure:"auto-offset-reset"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
	} `mapstructure:"kafka"`
}

func LoadConfig(fileName string) error {
	v := viper.New()

	v.SetConfigFile(fileName)
	v.SetConfigType("yaml")

	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("/opt/config")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return fmt.Errorf("error reading config file: %s", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return fmt.Errorf("unable to decode into struct: %v", err)
	}

	log.Printf("Config loaded successfully, %+v", config)

	// Start watching the config file for changes
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		if err := v.Unmarshal(&config); err != nil {
			log.Printf("Error decoding config after change, %v", err)
			return
		}

		log.Printf("Config reloaded successfully, %+v", config)
	})

	return nil
}

func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}
