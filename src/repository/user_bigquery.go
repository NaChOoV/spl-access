package repository

import (
	"context"
	"fmt"
	"os"
	"spl-access/src/dto"
	"strings"

	"cloud.google.com/go/bigquery"
)

type BigQueryUser struct {
	client *bigquery.Client
}

func NewBigQueryUser(client *bigquery.Client) *BigQueryUser {
	return &BigQueryUser{
		client: client,
	}
}

func (b *BigQueryUser) CheckUsers(ctx context.Context, users []*dto.UserDto, args ...any) error {
	var values []string
	for _, user := range users {
		value := fmt.Sprintf("('%s', '%s', '%s')",
			user.ExternalId,
			user.Run,
			user.FullName,
		)
		values = append(values, value)
	}

	bigQueryFile, err := os.ReadFile("src/repository/bigquery/user_upsert.sql")
	if err != nil {
		return err
	}

	// Create MERGE query
	mergeQuery := fmt.Sprintf(string(bigQueryFile), strings.Join(values, ", "))

	query := b.client.Query(mergeQuery)
	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to start merge job: %w", err)
	}

	// Wait for the job to complete
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for merge job: %w", err)
	}

	if status.Err() != nil {
		return fmt.Errorf("merge job failed: %w", status.Err())
	}

	return nil
}
