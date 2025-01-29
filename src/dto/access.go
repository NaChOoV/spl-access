package dto

import "time"

type AccessArrayDto struct {
	Data []AccessDto `json:"data" validate:"required,dive"`
}

type AccessDto struct {
	ExternalId string    `json:"externalId" validate:"required"`
	Location   string    `json:"location" validate:"required,oneof=102 104 105 106 107 108"`
	Run        string    `json:"run"  validate:"required"`
	FullName   string    `json:"fullName" validate:"required"`
	EntryAt    time.Time `json:"entryAt" validate:"required"`
	ExitAt     time.Time `json:"exitAt" validate:"omitempty"`
}
