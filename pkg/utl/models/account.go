package models

// Account represents account model
type Account struct {
	Base
	Name    string `json:"name"`
	Teams   []Team `json:"teams,omitempty"`
	OwnerID uint   `json:"owner_id"`
}
