package repository

import (
	"context"
	"spl-access/src/dto"
	"spl-access/src/model"
)

type AccessRepository interface {
	GetAccess(ctx context.Context, complete bool) ([]*model.Access, error)
	UpdateOrCreateAccess(ctx context.Context, access []*dto.AccessDto) error
}

type UserRepository interface {
	CheckUsers(ctx context.Context, users []*dto.UserDto, args ...any) error
}
