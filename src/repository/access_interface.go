package repository

import (
	"spl-access/src/dto"
	"spl-access/src/model"
)

type AccessRepository interface {
	GetAccess() (*[]model.Access, error)
	UpdateOrCreateAccess(access dto.AccessArrayDto) error
}
