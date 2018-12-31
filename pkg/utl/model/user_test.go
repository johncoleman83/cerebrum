package cerebrum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

func TestChangePassword(t *testing.T) {
	user := &cerebrum.User{
		FirstName: "TestGuy",
	}

	hashedPassword := "h4$h3D"

	user.ChangePassword(hashedPassword)
	if user.LastPasswordChange.IsZero() {
		t.Errorf("Last password change was not changed")
	}

	if user.Password != hashedPassword {
		t.Errorf("Password was not changed")

	}
}

func TestUpdateLastLogin(t *testing.T) {
	user := &cerebrum.User{
		FirstName: "TestGuy",
	}

	token := "helloWorld"

	user.UpdateLastLogin(token)
	if user.LastLogin.IsZero() {
		t.Errorf("Last login time was not changed")
	}

	if user.Token != token {
		t.Errorf("Tooken was not changed")

	}
}

func TestPaginationLimit(t *testing.T) {
	reqNegativeLimit := cerebrum.PaginationReq{Limit: -5, Page: 2}
	expected := &cerebrum.Pagination{Limit: 100, Offset: 200}
	assert.Equal(t, expected, reqNegativeLimit.NewPagination(), "negative limit should get set to default")

	reqMaxLimit := cerebrum.PaginationReq{Limit: 1001, Page: 2}
	expected.Limit, expected.Offset = 1000, 2000
	assert.Equal(t, expected, reqMaxLimit.NewPagination(), "beyond max limit should get set to default")

	reqTooBigLimit := cerebrum.PaginationReq{Limit: 9999999, Page: 2}
	expected.Limit, expected.Offset = 1000, 2000
	assert.Equal(t, expected, reqTooBigLimit.NewPagination(), "way beyond max limit should get set to default")

	reqNoChangeAllZeros := cerebrum.PaginationReq{Limit: 0, Page: 0}
	expected.Limit, expected.Offset = 100, 0
	assert.Equal(t, expected, reqNoChangeAllZeros.NewPagination(), "zeros should get set to default")

	reqNoChange := cerebrum.PaginationReq{Limit: 95, Page: 25}
	expected.Limit, expected.Offset = 95, 2375
	assert.Equal(t, expected, reqNoChange.NewPagination(), "some random offset and limit within the bounds should stay the same")
}
