package models

// Company represents company model
type Company struct {
	Base
	Name      string     `json:"name"`
	Locations []Location `json:"locations,omitempty"`
	OwnerID   uint       `json:"owner_id"`
}