package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port      int
	Database  string
	SentryDsn string
}

func loadConfig() *Config {
	cfg := &Config{
		Port:      1234,
		Database:  "myfile.db",
		SentryDsn: "https://565ab10db08448289861fe107cb0867b@o913183.ingest.sentry.io/6469771",
	}

	viper.SetConfigName("front-test")
	viper.SetEnvPrefix("ft")

	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE")

	//Flags
	viper.AutomaticEnv()

	viper.ReadInConfig()

	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Println("cannot unmarshal config: %s", err)
	}

	return cfg
}
