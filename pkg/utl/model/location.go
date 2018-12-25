package cerebrum

import (
	"github.com/jinzhu/gorm"
)

// Location represents company location model
type Location struct {
	gorm.Model
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Address string `json:"address"`

	CompanyID int `json:"company_id"`
}
