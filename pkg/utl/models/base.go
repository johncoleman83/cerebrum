package models

import "time"

// Base contains common fields for all tables
type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// ListQuery holds account/team data used for list db queries
type ListQuery struct {
	Query string
	ID    uint
}
