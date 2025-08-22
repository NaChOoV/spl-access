package service

import (
	"net/http"
	"spl-access/src/config"
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
