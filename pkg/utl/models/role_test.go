package models_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/stretchr/testify/assert"
)

func TestNewRoleFromAccessLevelUint(t *testing.T) {
	correctMessage := "should return the correct role for input AccessLevel uint"
	errorMessage := "should return error when input is not a known AccessLevel"

	actual, _ := models.NewRoleFromAccessLevelUint(100)
	assert.Equal(t, uint(1), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromAccessLevelUint(110)
	assert.Equal(t, uint(2), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromAccessLevelUint(120)
	assert.Equal(t, uint(3), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromAccessLevelUint(130)
	assert.Equal(t, uint(4), actual.ID, correctMessage)

	actual, _ = models.NewRoleFromAccessLevelUint(200)
	assert.Equal(t, uint(5), actual.ID, correctMessage)

	_, err := models.NewRoleFromAccessLevelUint(0)
	assert.Equal(t, "unknown accessLevel id", err.Error(), errorMessage)

	_, err = models.NewRoleFromAccessLevelUint(10)
	assert.Equal(t, "unknown accessLevel id", err.Error(), errorMessage)

	_, err = models.NewRoleFromAccessLevelUint(300)
	assert.Equal(t, "unknown accessLevel id", err.Error(), errorMessage)

	_, err = models.NewRoleFromAccessLevelUint(2000)
	assert.Equal(t, "unknown accessLevel id", err.Error(), errorMessage)
}
