package models

// Location represents company location model
type Location struct {
	Base
	Name      string `json:"name"`
	Address   string `json:"address"`
	CompanyID uint   `json:"company_id"`
}
