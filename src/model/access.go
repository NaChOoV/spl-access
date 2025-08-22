package model

import (
	"time"
)

type Access struct {
	ExternalId string     `json:"externalId"`
	Run        string     `json:"run"`
	FullName   string     `json:"fullName"`
	Location   string     `json:"location"`
	EntryAt    time.Time  `json:"entryAt"`
	ExitAt     *time.Time `json:"exitAt"`
}
