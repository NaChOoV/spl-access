package config

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type EnvironmentConfig struct {
	AuthString string `env:"AUTH_STRING,required"`
	DebugMode  bool   `env:"DEBUG_MODE"`
	DbHost     string `env:"DB_HOST,required"`
	DbPort     string `env:"DB_PORT,required"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`
}

var envConfig *EnvironmentConfig

func NewEnviromentConfig(lc fx.Lifecycle) *EnvironmentConfig {
	envConfig = &EnvironmentConfig{}

	var err error
	isLocalhost := os.Getenv("IS_LOCALHOST")
	if isLocalhost != "true" {
		err = godotenv.Load()
	}

	if err != nil {
		panic(err)
	}
	// AuthString
	envConfig.AuthString = os.Getenv("AUTH_STRING")

	// Database
	envConfig.DbHost = os.Getenv("DB_HOST")
	envConfig.DbPort = os.Getenv("DB_PORT")
	envConfig.DbUser = os.Getenv("DB_USER")
	envConfig.DbPassword = os.Getenv("DB_PASSWORD")
	envConfig.DbName = os.Getenv("DB_NAME")

	// DebugMode
	if os.Getenv("DEBUG_MODE") == "true" {
		envConfig.DebugMode = true
	} else {
		envConfig.DebugMode = false
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if err != nil {
				return err
			}
			return nil
		},
	})

	printEnvironmentConfig(*envConfig)
	return envConfig
}

func printEnvironmentConfig(config EnvironmentConfig) {
	v := reflect.ValueOf(config)
	typeOfConfig := v.Type()

	fmt.Println("EnvironmentConfig:")
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("  %s: %v\n", typeOfConfig.Field(i).Name, v.Field(i).Interface())
	}
}

func NewContextBackground() *context.Context {
	ctx := context.Background()
	return &ctx
}
