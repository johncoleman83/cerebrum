package models

import (
	"errors"
)

// AccessRole represents access role type
type AccessRole uint

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 100

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// AccountAdminRole can edit account specific things
	AccountAdminRole AccessRole = 120

	// TeamAdminRole can edit team specific things
	TeamAdminRole AccessRole = 130

	// UserRole is a standard user
	UserRole AccessRole = 200
)

// ValidRoles contains all valid roles mapped to their ID
var ValidRoles = map[uint]Role{
	100: Role{ID: 1, AccessLevel: SuperAdminRole, Name: "SUPER_ADMIN"},
	110: Role{ID: 2, AccessLevel: AdminRole, Name: "ADMIN"},
	120: Role{ID: 3, AccessLevel: AccountAdminRole, Name: "ACCOUNT_ADMIN"},
	130: Role{ID: 4, AccessLevel: TeamAdminRole, Name: "TEAM_ADMIN"},
	200: Role{ID: 5, AccessLevel: UserRole, Name: "USER_ADMIN"},
}

// Role model
type Role struct {
	ID          uint       `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}

// NewRoleFromAccessLevelUint contains all valid roles
func NewRoleFromAccessLevelUint(accessLevel uint) (*Role, error) {
	role, ok := ValidRoles[accessLevel]
	if !ok {
		return nil, errors.New("unknown accessLevel id")
	}
	return &role, nil
}
