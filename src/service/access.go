package service

import (
	"fmt"
	"spl-access/src/config"
	"spl-access/src/dto"
	"spl-access/src/helpers"
	"spl-access/src/model"
	"spl-access/src/repository"
	"spl-access/src/websocket"
)

type AccessService struct {
	accessRepository    *repository.AccessRepository
	websocketController *websocket.WebsocketController
	config              *config.EnvironmentConfig
	access              *[]model.Access
}

func NewAccessService(
	accessRepository *repository.AccessRepository,
	websocketController *websocket.WebsocketController,
	config *config.EnvironmentConfig,
) *AccessService {
	return &AccessService{
		accessRepository:    accessRepository,
		websocketController: websocketController,
		config:              config,
		access:              &[]model.Access{},
	}
}

func (a *AccessService) UpdateOrCreateAccess(access dto.AccessArrayDto) error {
	cleanedAccess := a.RemoveDuplicates(access)
	err := a.accessRepository.UpdateOrCreateAccess(cleanedAccess)
	if err != nil {
		return err
	}

	go a.UpdateAccess()
	return nil
}

func (a *AccessService) UpdateAccess() {

	access, err := a.accessRepository.GetAccess()
	if err != nil {
		fmt.Println(err)
		return
	}

	a.access = access
	obfuscateAccess := helpers.MaskAccessData(access)

	a.websocketController.BroadcastMessage(obfuscateAccess)
}

func (a *AccessService) GetTodayAccess() *[]model.Access {
	return a.access
}

func (a *AccessService) RemoveDuplicates(access dto.AccessArrayDto) dto.AccessArrayDto {
	seen := make(map[string]bool)
	uniqueData := []dto.AccessDto{}

	for _, entry := range access.Data {
		key := entry.Run + entry.EntryAt.String()
		if !seen[key] {
			seen[key] = true
			uniqueData = append(uniqueData, entry)
		}
	}

	return dto.AccessArrayDto{Data: uniqueData}
}
