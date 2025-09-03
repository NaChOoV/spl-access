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
	Port                   string `env:"PORT,default=8000"`
	AuthString             string `env:"AUTH_STRING,required"`
	SourceBaseURL          string `env:"SOURCE_BASE_URL,required"`
	SourceAuthString       string `env:"SOURCE_AUTH_STRING,required"`
	NotificationBaseUrl    string `env:"NOTIFICATION_BASE_URL,required"`
	NotificationAuthString string `env:"NOTIFICATION_AUTH_STRING,required"`
	DebugMode              bool   `env:"DEBUG_MODE"`
	DbHost                 string `env:"DB_HOST"`
	DbPort                 string `env:"DB_PORT"`
	DbUser                 string `env:"DB_USER"`
	DbPassword             string `env:"DB_PASSWORD"`
	DbName                 string `env:"DB_NAME"`
	DbSSLMode              string `env:"DB_SSL_MODE"`
	BigQueryProjectID      string `env:"BIG_QUERY_PROJECT_ID,required"`

	Zone string `env:"ZONE"`
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

	// Port
	envConfig.Port = os.Getenv("PORT")
	if envConfig.Port == "" {
		envConfig.Port = "8000"
	}

	// AuthString
	envConfig.AuthString = os.Getenv("AUTH_STRING")

	// Source service
	envConfig.SourceBaseURL = os.Getenv("SOURCE_BASE_URL")
	envConfig.SourceAuthString = os.Getenv("SOURCE_AUTH_STRING")

	// Notification Service
	envConfig.NotificationBaseUrl = os.Getenv("NOTIFICATION_BASE_URL")
	envConfig.NotificationAuthString = os.Getenv("NOTIFICATION_AUTH_STRING")

	// Database
	envConfig.DbHost = os.Getenv("DB_HOST")
	envConfig.DbPort = os.Getenv("DB_PORT")
	envConfig.DbUser = os.Getenv("DB_USER")
	envConfig.DbPassword = os.Getenv("DB_PASSWORD")
	envConfig.DbName = os.Getenv("DB_NAME")
	envConfig.DbSSLMode = os.Getenv("DB_SSL_MODE")

	// BigQueryProjectID
	envConfig.BigQueryProjectID = os.Getenv("BIG_QUERY_PROJECT_ID")

	// DebugMode
	if os.Getenv("DEBUG_MODE") == "true" {
		envConfig.DebugMode = true
	} else {
		envConfig.DebugMode = false
	}

	// Zone
	envConfig.Zone = os.Getenv("ZONE")
	if envConfig.Zone == "" {
		envConfig.Zone = "GMT-3"
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
