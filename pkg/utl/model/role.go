package cerebrum

// AccessRole represents access role type
type AccessRole uint

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 100

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// CompanyAdminRole can edit company specific things
	CompanyAdminRole AccessRole = 120

	// LocationAdminRole can edit location specific things
	LocationAdminRole AccessRole = 130

	// UserRole is a standard user
	UserRole AccessRole = 200
)

// ValidRoles contains all valid roles mapped to their ID
var ValidRoles = map[uint]uint{
	100: 1,
	110: 2,
	120: 3,
	130: 4,
	200: 5,
}

// Role model
type Role struct {
	ID          uint       `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}

// AccessLevelToID contains all valid roles
func AccessLevelToID(accessLevel uint) uint {
	id, ok := ValidRoles[accessLevel]
	if !ok {
		return 0
	}
	return id
}
