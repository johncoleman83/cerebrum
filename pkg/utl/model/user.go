package cerebrum

import (
	"time"
)

// User represents user domain model
type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`

	Mobile  string `json:"mobile,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`

	Active bool `json:"active"`

	LastLogin          time.Time `json:"last_login,omitempty"`
	LastPasswordChange time.Time `json:"last_password_change,omitempty"`

	Token string `json:"-"`

	Role *Role `json:"role,omitempty"`

	RoleID     AccessRole `json:"-"`
	CompanyID  uint        `json:"company_id"`
	LocationID uint        `json:"location_id"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID         uint
	CompanyID  uint
	LocationID uint
	Username   string
	Email      string
	Role       AccessRole
}

// ChangePassword updates user's password related fields
func (u *User) ChangePassword(hash string) {
	u.Password = hash
	u.LastPasswordChange = time.Now()
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin(token string) {
	u.Token = token
	u.LastLogin = time.Now()
}
