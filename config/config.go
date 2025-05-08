package config

import (
    "github.com/spf13/viper"
    "log"
)

type H3Config struct {
    ResolutionMin int `mapstructure:"resolution_min"`
    ResolutionMax int `mapstructure:"resolution_max"`
}

type Config struct {
    H3 H3Config `mapstructure:"h3"`
}

var AppConfig Config

func LoadConfig() {
    viper.SetConfigName("config")      
    viper.SetConfigType("yaml")        
    viper.AddConfigPath("./config")    

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    if err := viper.Unmarshal(&AppConfig); err != nil {
        log.Fatalf("Error unmarshaling config: %v", err)
    }
}
