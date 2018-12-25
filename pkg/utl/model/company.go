package cerebrum

import (
	"github.com/jinzhu/gorm"
)

// Company represents company model
type Company struct {
	gorm.Model
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	Locations []Location `json:"locations,omitempty"`
	Owner     User       `json:"owner"`
}
