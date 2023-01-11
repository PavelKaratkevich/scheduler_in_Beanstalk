package util

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
}

func ConnectDB(envVarsPath string) (*sql.DB, error) {

	config, err := LoadEnvVars(envVarsPath)
	if err != nil {
		log.Fatal("Error while loading env variables:", err.Error())
	}

	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error while opening db: ", err.Error())
		return nil, err
	}

	return db, nil
}

func LoadEnvVars(path string) (config Config, err error) {
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("env") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)  // path to look for the config file in

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.AutomaticEnv() // allows to overwrite env variable from command line

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Unable to load ennvironment variables: ", err.Error())
		return
	}

	return config, err
}
