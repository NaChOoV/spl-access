package model

import "time"

type Access struct {
	Run      string    `json:"run"`
	FullName string    `json:"fullName"`
	Location string    `json:"location"`
	EntryAt  time.Time `json:"entryAt"`
}
