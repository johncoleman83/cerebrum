package models

// Team represents account team model
type Team struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountID   uint   `json:"account_id"`
}
