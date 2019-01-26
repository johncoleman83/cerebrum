package models_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/stretchr/testify/assert"
)

func TestNewRoleFromRoleID(t *testing.T) {
	correctMessage := "should return the correct role for input AccessLevel uint"
	errorMessage := "should return error when input is not a known AccessLevel"

	actual, _ := models.NewRoleFromRoleID(1)
	assert.Equal(t, uint(1), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromRoleID(2)
	assert.Equal(t, uint(2), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromRoleID(3)
	assert.Equal(t, uint(3), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromRoleID(4)
	assert.Equal(t, uint(4), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromRoleID(5)
	assert.Equal(t, uint(5), actual.ID, correctMessage)

	_, err := models.NewRoleFromRoleID(0)
	assert.Equal(t, "unknown role id", err.Error(), errorMessage)

	_, err = models.NewRoleFromRoleID(6)
	assert.Equal(t, "unknown role id", err.Error(), errorMessage)

	_, err = models.NewRoleFromRoleID(30)
	assert.Equal(t, "unknown role id", err.Error(), errorMessage)

	_, err = models.NewRoleFromRoleID(2000)
	assert.Equal(t, "unknown role id", err.Error(), errorMessage)
}
