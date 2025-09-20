package db

import (
	"context"
	"fmt"
	"spl-access/ent"
	"spl-access/src/config"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func CreatePostgresConnection(lc fx.Lifecycle, envConfig *config.EnvironmentConfig, ctx *context.Context) *ent.Client {
	options := []ent.Option{}
	if envConfig.DebugMode {
		options = append(options, ent.Debug())
	}
	conn, connError := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		envConfig.DbHost,
		envConfig.DbPort,
		envConfig.DbUser,
		envConfig.DbName,
		envConfig.DbPassword,
		envConfig.DbSSLMode),
		options...)

	// Run the auto migration tool.
	migrationError := conn.Schema.Create(*ctx)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if connError != nil {
				fmt.Printf("Failed opening connection to postgres: %v", connError)
				return connError
			}
			if migrationError != nil {
				fmt.Printf("Failed to execute the migration tool: %v", migrationError)
				return migrationError
			}
			return nil
		},
		OnStop: func(context.Context) error {
			err := conn.Close()
			if err != nil {
				fmt.Printf("Error closing the database connection: %v", err)
				return err
			}
			fmt.Println("Database connection closed.")
			return nil
		},
	})

	return conn
}
