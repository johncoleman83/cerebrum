package transport

import (
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*cerebrum.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []cerebrum.User `json:"users"`
		Page  int          `json:"page"`
	}
}
