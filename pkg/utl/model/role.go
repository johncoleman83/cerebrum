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

// ValidRoles contains all valid rols
var ValidRoles = map[AccessRole]bool{
	SuperAdminRole:    true,
	AdminRole:         true,
	CompanyAdminRole:  true,
	LocationAdminRole: true,
	UserRole:          true,
}

// Role model
type Role struct {
	ID          AccessRole `json:"id" gorm:"foreignkey:RoleID;association_foreignkey:ID;"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}
