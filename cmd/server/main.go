package main

import (
	"fmt"
	"github.com/inasknh/simple-poke-app/internal/config"
	"github.com/spf13/viper"
	"log"
)

func main() {
	fmt.Println("hello")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read config")
	}

	var configuration config.Configurations
	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Couldn't unmarshal configuration")
	}

}
