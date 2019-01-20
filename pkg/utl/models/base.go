package models

import (
	"github.com/jinzhu/gorm"
)

// Base contains common fields for all tables
type Base struct {
	gorm.Model
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    uint
}
