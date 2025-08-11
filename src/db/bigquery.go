package db

import (
	"context"
	"fmt"
	"spl-access/src/config"

	"cloud.google.com/go/bigquery"
	"go.uber.org/fx"
)

func CreateBigQueryConnection(lc fx.Lifecycle, envConfig *config.EnvironmentConfig, ctx *context.Context) (*bigquery.Client, error) {
	client, err := bigquery.NewClient(*ctx, envConfig.BigQueryProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}

	if err := bigQuerySchemaMigration(ctx, client, envConfig); err != nil {
		return nil, fmt.Errorf("failed to run BigQuery schema migration: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

func bigQuerySchemaMigration(ctx *context.Context, client *bigquery.Client, envConfig *config.EnvironmentConfig) error {
	// Implement schema migration logic here
	// This is a placeholder function; actual implementation will depend on your requirements

	datasetID := "access" // Replace with your actual dataset ID
	tableID := "access"

	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)

	// Check if table exists
	_, err := table.Metadata(*ctx)
	if err != nil {
		// Table doesn't exist, create it
		schema := bigquery.Schema{
			{Name: "external_id", Type: bigquery.StringFieldType, Required: true},
			{Name: "run", Type: bigquery.StringFieldType, Required: true},
			{Name: "location", Type: bigquery.IntegerFieldType, Required: true},
			{Name: "entry_at", Type: bigquery.TimestampFieldType, Required: true},
			{Name: "exit_at", Type: bigquery.TimestampFieldType, Required: false},
		}

		tableMetadata := &bigquery.TableMetadata{
			Schema: schema,
		}

		if err := table.Create(*ctx, tableMetadata); err != nil {
			return fmt.Errorf("failed to create access table: %w", err)
		}

		fmt.Println("access table created successfully.")
	}

	tableID = "today_access"
	table = dataset.Table(tableID)
	// Check if user table exists
	_, err = table.Metadata(*ctx)
	if err != nil {
		// User table doesn't exist, create it
		schema := bigquery.Schema{
			{Name: "external_id", Type: bigquery.StringFieldType, Required: true},
			{Name: "run", Type: bigquery.StringFieldType, Required: true},
			{Name: "location", Type: bigquery.IntegerFieldType, Required: true},
			{Name: "entry_at", Type: bigquery.TimestampFieldType, Required: true},
			{Name: "exit_at", Type: bigquery.TimestampFieldType, Required: false},
		}

		tableMetadata := &bigquery.TableMetadata{
			Schema: schema,
		}

		if err := table.Create(*ctx, tableMetadata); err != nil {
			return fmt.Errorf("failed to create user table: %w", err)
		}

		fmt.Println("today_access table created successfully.")
	}

	tableID = "user"
	table = dataset.Table(tableID)
	// Check if user table exists
	_, err = table.Metadata(*ctx)
	if err != nil {
		// User table doesn't exist, create it
		schema := bigquery.Schema{
			{Name: "run", Type: bigquery.StringFieldType, Required: true},
			{Name: "external_id", Type: bigquery.StringFieldType, Required: true},
			{Name: "full_name", Type: bigquery.StringFieldType, Required: true},
		}

		tableMetadata := &bigquery.TableMetadata{
			Schema: schema,
		}

		if err := table.Create(*ctx, tableMetadata); err != nil {
			return fmt.Errorf("failed to create user table: %w", err)
		}

		fmt.Println("user table created successfully.")
	}

	return nil
}
