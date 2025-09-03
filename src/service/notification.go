package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spl-access/src/config"
	"spl-access/src/model"
	"time"
)

type NotificationService struct {
	httpClient *http.Client
	config     *config.EnvironmentConfig
}

func NewNotificationService(config *config.EnvironmentConfig) *NotificationService {
	return &NotificationService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

func (ns *NotificationService) SendAccess(ctx context.Context, acceses []*model.Access) error {
	// Implementation for sending a notification

	url := fmt.Sprintf("%s/access", ns.config.NotificationBaseUrl)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", ns.config.NotificationAuthString)
	req.Header.Set("Content-Type", "application/json")
	jsonBody, err := json.Marshal(acceses)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	req.ContentLength = int64(len(jsonBody))

	resp, err := ns.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification service failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}
