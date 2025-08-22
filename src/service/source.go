package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spl-access/src/config"
	"spl-access/src/dto"
	"spl-access/src/types"
	"time"
)

type AccessResponse struct {
	Message string `json:"message"`
	Data    struct {
		Records []types.AccessRecord `json:"records"`
	} `json:"data"`
}

type SourceService struct {
	httpClient *http.Client
	config     *config.EnvironmentConfig
}

func NewSourceService(config *config.EnvironmentConfig) *SourceService {
	return &SourceService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

func (s *SourceService) GetAccess() ([]*dto.AccessDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.SourceBaseURL+"/access", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("X-Auth-String", s.config.SourceAuthString)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var response AccessResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert AccessRecord to model.Access
	var accessList []*dto.AccessDto
	for _, record := range response.Data.Records {
		// Parse timestamps
		entryAt, err := time.Parse(time.RFC3339, record.EntryAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse entryAt timestamp: %w", err)
		}

		access := &dto.AccessDto{
			ExternalId: fmt.Sprintf("%d", record.ExternalId),
			Location:   record.Location,
			Run:        record.Run,
			FullName:   record.FullName,
			EntryAt:    entryAt,
		}

		// Parse exitAt if present and not empty
		if record.ExitAt != "" {
			exitAt, err := time.Parse(time.RFC3339, record.ExitAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse exitAt timestamp: %w", err)
			}
			access.ExitAt = &exitAt
		}
		// If ExitAt is empty or null, access.ExitAt remains nil

		accessList = append(accessList, access)
	}

	return accessList, nil
}
