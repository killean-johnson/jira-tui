package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
    Email string
    APIToken string
    JiraURL string
    Board []LayoutStruct
    Issue []LayoutStruct
}

type LayoutStruct struct {
    View string
    Keys []Keybinding
}

type Keybinding struct {
    Name string
    Key string
    Description string
}

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func (kb *Config) LoadConfig() error {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    viper.SetConfigType("json")

    viper.AutomaticEnv()

    err := viper.ReadInConfig()
    if err != nil {
        return err
    }

    err = viper.Unmarshal(&kb)
    if err != nil {
        return err
    }

    return nil
}
