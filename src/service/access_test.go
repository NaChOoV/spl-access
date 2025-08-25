package service

import (
	"context"
	"errors"
	"spl-access/src/config"
	"spl-access/src/dto"
	"spl-access/src/helpers"
	"spl-access/src/model"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock AccessRepository ---
type MockAccessRepository struct {
	mock.Mock
}

func (m *MockAccessRepository) GetAccess(ctx context.Context) ([]*model.Access, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Access), args.Error(1)

}

func (m *MockAccessRepository) UpdateOrCreateAccess(ctx context.Context, access []*dto.AccessDto) error {
	args := m.Called(ctx, access)
	return args.Error(0)
}

// --- Mock WebSocketController ---
type MockWebsocketController struct {
	mock.Mock
}

func (m *MockWebsocketController) BroadcastMessage(data any) {
	m.Called(data)
}
func (m *MockWebsocketController) Upgrade(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestUpdateAccess(t *testing.T) {
	mockConfig := &config.EnvironmentConfig{}

	now := time.Now()
	originalAccessData := []*model.Access{
		{Run: "11111111-1", FullName: "Test User One", Location: "Location1", EntryAt: now},
		{Run: "22222222-2", FullName: "Test User Two", Location: "Location2", EntryAt: now.Add(-time.Hour)},
	}
	emptyAccessData := []*model.Access{}

	unexpectedError := errors.New("unexpected error")

	expectedMaskedOriginalData := helpers.MaskAccessData(originalAccessData)
	expectedMaskedEmptyData := helpers.MaskAccessData(emptyAccessData)
	preExistingAccessState := []*model.Access{{Run: "PRE-EXISTING-000", FullName: "Old Data"}}

	testCases := []struct {
		name                  string
		setupMocks            func(mockRepo *MockAccessRepository, mockWS *MockWebsocketController)
		initialServiceAccess  []*model.Access
		expectedServiceAccess []*model.Access
	}{
		{
			setupMocks: func(mockRepo *MockAccessRepository, mockWS *MockWebsocketController) {
				dataCopy := make([]*model.Access, len(originalAccessData))
				copy(dataCopy, originalAccessData)
				mockRepo.On("GetAccess", mock.Anything).Return(dataCopy, nil).Once()
				mockWS.On("BroadcastMessage", expectedMaskedOriginalData).Return().Once()
			},
			initialServiceAccess:  []*model.Access{},
			expectedServiceAccess: originalAccessData,
		},
		{
			setupMocks: func(mockRepo *MockAccessRepository, mockWS *MockWebsocketController) {
				mockRepo.On("GetAccess", mock.Anything).Return(nil, unexpectedError).Once()
			},
			initialServiceAccess:  preExistingAccessState,
			expectedServiceAccess: preExistingAccessState,
		},
		{
			setupMocks: func(mockRepo *MockAccessRepository, mockWS *MockWebsocketController) {
				mockRepo.On("GetAccess", mock.Anything).Return([]*model.Access{}, nil).Once()
				mockWS.On("BroadcastMessage", expectedMaskedEmptyData).Return().Once()
			},
			initialServiceAccess:  preExistingAccessState,
			expectedServiceAccess: emptyAccessData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			mockRepo := new(MockAccessRepository)
			mockWS := new(MockWebsocketController)

			tc.setupMocks(mockRepo, mockWS)

			serviceAccessCopy := make([]*model.Access, len(tc.initialServiceAccess))
			copy(serviceAccessCopy, tc.initialServiceAccess)

			service := &AccessService{
				accessRepository:    mockRepo,
				websocketController: mockWS,
				config:              mockConfig,
				access:              serviceAccessCopy,
			}

			service.UpdateAccess(context.Background())

			mockRepo.AssertExpectations(t)
			mockWS.AssertExpectations(t)

			assert.Equal(tc.expectedServiceAccess, service.access, "service.access state mismatch")
		})
	}
}
