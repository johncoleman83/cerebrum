package models_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/stretchr/testify/assert"
)

/*
	Test Cases
	100: 1,
	110: 2,
	120: 3,
	130: 4,
	200: 5,
*/
func TestAccessLevelToID(t *testing.T) {
	correctMessage := "should return the correct uint ID for input AccessLevel uint"
	errorMessage := "should return 0 when input is not a known AccessLevel"

	assert.Equal(t, uint(1), models.AccessLevelToID(100), correctMessage)
	assert.Equal(t, uint(2), models.AccessLevelToID(110), correctMessage)
	assert.Equal(t, uint(3), models.AccessLevelToID(120), correctMessage)
	assert.Equal(t, uint(4), models.AccessLevelToID(130), correctMessage)
	assert.Equal(t, uint(5), models.AccessLevelToID(200), correctMessage)
	assert.Equal(t, uint(0), models.AccessLevelToID(0), errorMessage)
	assert.Equal(t, uint(0), models.AccessLevelToID(10), errorMessage)
	assert.Equal(t, uint(0), models.AccessLevelToID(300), errorMessage)
	assert.Equal(t, uint(0), models.AccessLevelToID(2000), errorMessage)
}
