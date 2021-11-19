package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	User   UserConfig `mapstructure:"user"`
	Server string     `mapstructure:"server"`
}

type UserConfig struct {
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"`
}

// Try to load config from ~/.config/tira/tira_config.json
func LoadConfig() {
	//build path to config file
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "tira")
	configFile := filepath.Join(configPath, "tira_config.json")

	viper.SetConfigName("tira_config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	//if config file doenst exist, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		//create folder structure to file
		if err := os.MkdirAll(filepath.Dir(configFile), 0770); err != nil {
			log.Fatal("Error creating config file")
		}

		//create file
		createDefaultConfig(configFile)
	}

	// read in the config data
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Error loading config:", err)
		}
	}
}

// separate function to create default config file
// makes managing defaults easier
func createDefaultConfig(configFile string) {
	os.Create(configFile)
	viper.SetDefault("User.Username", "")
	viper.SetDefault("User.Token", "")
	viper.SetDefault("Server", "")

	viper.WriteConfig()
}

// Get the config
func GetConfig() (config Config) {
	err := viper.Unmarshal(&config)

	if err != nil {
		panic("Error parsing config")
	}
	return
}
