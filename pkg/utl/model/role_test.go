package cerebrum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
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

	assert.Equal(t, uint(1), cerebrum.AccessLevelToID(100), correctMessage)
	assert.Equal(t, uint(2), cerebrum.AccessLevelToID(110), correctMessage)
	assert.Equal(t, uint(3), cerebrum.AccessLevelToID(120), correctMessage)
	assert.Equal(t, uint(4), cerebrum.AccessLevelToID(130), correctMessage)
	assert.Equal(t, uint(5), cerebrum.AccessLevelToID(200), correctMessage)
	assert.Equal(t, uint(0), cerebrum.AccessLevelToID(0), errorMessage)
	assert.Equal(t, uint(0), cerebrum.AccessLevelToID(10), errorMessage)
	assert.Equal(t, uint(0), cerebrum.AccessLevelToID(300), errorMessage)
	assert.Equal(t, uint(0), cerebrum.AccessLevelToID(2000), errorMessage)
}
