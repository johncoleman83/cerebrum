package support_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

func TestNewRandomString(t *testing.T) {
	t.Run("test new random string should always produce new random strings of correct length", func(t *testing.T) {
		str1 := support.NewRandomString(-1)
		assert.Equal(t, "", str1, "input of less than 1 should not break everything")
		assert.Equal(t, 0, len(str1), "the length should be the same as the input integer provided")

		str2 := support.NewRandomString(0)
		assert.Equal(t, "", str2, "input of less than 1 should not break everything")

		strLong := support.NewRandomString(999)
		assert.Equal(t, 999, len(strLong), "the length should be the same as the input integer provided")

		strRand1 := support.NewRandomString(25)
		assert.Equal(t, 25, len(strRand1), "the length should be the same as the input integer provided")

		strRand2 := support.NewRandomString(25)
		assert.Equal(t, 25, len(strRand2), "the length should be the same as the input integer provided")

		assert.NotEqual(t, strRand1, strRand2, "strings should never be the same")
	})
}
