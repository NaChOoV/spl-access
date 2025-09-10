package service

import (
	"context"
	"fmt"
	"spl-access/src/config"
	"spl-access/src/dto"
	"spl-access/src/helpers"
	"spl-access/src/model"
	"spl-access/src/repository"
	"spl-access/src/websocket"
	"sync"
)

type AccessService struct {
	accessRepository    repository.AccessRepository
	userRepository      repository.UserRepository
	sourceService       *SourceService
	notificationService *NotificationService
	websocketController websocket.AccessWb
	config              *config.EnvironmentConfig
	access              []*model.Access
	accessComplete      []*model.Access
}

func NewAccessService(
	accessRepository repository.AccessRepository,
	userRepository repository.UserRepository,
	sourceService *SourceService,
	notificationService *NotificationService,
	websocketController websocket.AccessWb,
	config *config.EnvironmentConfig,
) *AccessService {
	return &AccessService{
		accessRepository:    accessRepository,
		userRepository:      userRepository,
		sourceService:       sourceService,
		notificationService: notificationService,
		websocketController: websocketController,
		config:              config,
		access:              []*model.Access{},
		accessComplete:      []*model.Access{},
	}
}

func (a *AccessService) UpdateOrCreateAccess(ctx context.Context, access []*dto.AccessDto) error {
	cleanedAccess := helpers.RemoveDuplicatesGeneric(access, func(entry *dto.AccessDto) string {
		return entry.ExternalId + entry.Location + entry.EntryAt.String()
	})

	var accessErr, userErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		accessErr = a.accessRepository.UpdateOrCreateAccess(ctx, cleanedAccess)
	}()

	// CheckUsers in goroutine
	go func() {
		defer wg.Done()
		// Check Users
		var userMap = make(map[string]dto.UserDto)
		for _, accessItem := range cleanedAccess {
			if _, exists := userMap[accessItem.ExternalId]; !exists {
				userMap[accessItem.ExternalId] = dto.UserDto{
					Run:        accessItem.Run,
					ExternalId: accessItem.ExternalId,
					FullName:   accessItem.FullName,
				}
			}
		}

		var users []*dto.UserDto
		for _, user := range userMap {
			users = append(users, &user)
		}

		userErr = a.userRepository.CheckUsers(ctx, users)
	}()

	// Wait for both operations to complete
	wg.Wait()

	// Check errors
	if accessErr != nil {
		return accessErr
	}
	if userErr != nil {
		return userErr
	}

	go a.UpdateAccess(ctx)
	return nil
}

func (a *AccessService) UpdateAccess(ctx context.Context) {
	type result struct {
		access            []*model.Access
		accessComplete    []*model.Access
		accessErr         error
		accessCompleteErr error
	}

	var res result
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	go func() {
		defer wg.Done()
		accessResponse, err := a.accessRepository.GetAccess(ctx, false)
		mu.Lock()
		res.access = accessResponse
		res.accessErr = err
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		accessCompleteResponse, err := a.accessRepository.GetAccess(ctx, true)
		mu.Lock()
		res.accessComplete = accessCompleteResponse
		res.accessCompleteErr = err
		mu.Unlock()
	}()

	wg.Wait()

	if res.accessErr != nil {
		fmt.Println(res.accessErr)
		return
	}
	if res.accessCompleteErr != nil {
		fmt.Println(res.accessCompleteErr)
		return
	}

	a.access = res.access
	a.accessComplete = res.accessComplete

	obfuscateAccess := helpers.MaskAccessData(res.access)
	a.websocketController.BroadcastMessage(obfuscateAccess)
}

func (a *AccessService) GetTodayAccess() []*model.Access {
	return a.access
}

func (a *AccessService) GetTodayAccessComplete() []*model.Access {
	return a.accessComplete
}

func (a *AccessService) SyncAccess() error {
	ctx := context.Background()
	helpers.CheckSleepTime(a.config.Zone)

	// Get access from source service
	access, err := a.sourceService.GetAccess()
	if err != nil {
		return err
	}

	if len(access) == 0 {
		fmt.Println("No access data received from source service.")
		return nil
	}

	// Call UpdateOrCreateAccess with the access response.
	err = a.UpdateOrCreateAccess(ctx, access)
	if err != nil {
		return err
	}

	return nil
}
