package repository

import (
	"context"
	"fmt"
	"os"
	"spl-access/src/dto"
	"spl-access/src/model"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQueryAccess struct {
	client *bigquery.Client
}

func NewBigQueryAccess(client *bigquery.Client, userRepository UserRepository) *BigQueryAccess {
	return &BigQueryAccess{
		client: client,
	}
}

func (b *BigQueryAccess) UpdateOrCreateAccess(ctx context.Context, access *[]dto.AccessDto) error {
	var values []string
	for _, accessDto := range *access {
		exitAtValue := "TIMESTAMP(NULL)"
		if !accessDto.ExitAt.IsZero() {
			exitAtValue = fmt.Sprintf("TIMESTAMP('%s')", accessDto.ExitAt.Format(time.RFC3339))
		}

		value := fmt.Sprintf("('%s', '%s', CAST('%s' AS INT64), TIMESTAMP('%s'), %s)",
			accessDto.ExternalId,
			accessDto.Run,
			accessDto.Location,
			accessDto.EntryAt.Format(time.RFC3339),
			exitAtValue,
		)
		values = append(values, value)
	}

	bigQueryFile, err := os.ReadFile("src/repository/bigquery/today_access_upsert.sql")
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

func (b *BigQueryAccess) GetAccess(ctx context.Context) (*[]model.Access, error) {
	bigQueryFile, err := os.ReadFile("src/repository/bigquery/access_get.sql")
	if err != nil {
		return nil, err
	}

	query := b.client.Query(string(bigQueryFile))
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read from today_access table: %w", err)
	}

	var accessList []model.Access
	for {
		var row struct {
			Run      string    `bigquery:"run"`
			FullName string    `bigquery:"full_name"`
			Location int64     `bigquery:"location"`
			EntryAt  time.Time `bigquery:"entry_at"`
			ExitAt   time.Time `bigquery:"exit_at"`
		}

		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over results: %w", err)
		}

		access := model.Access{
			Run:      row.Run,
			FullName: row.FullName,
			Location: fmt.Sprintf("%d", row.Location),
			EntryAt:  row.EntryAt,
			ExitAt:   row.ExitAt,
		}
		accessList = append(accessList, access)
	}

	return &accessList, nil
}
