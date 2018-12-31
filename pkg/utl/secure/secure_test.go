package secure_test

import (
	"crypto/sha1"
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/secure"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	cases := []struct {
		name     string
		pass     string
		inputs   []string
		expected bool
	}{
		{
			name:     "Insecure password",
			pass:     "notSec",
			expected: false,
		},
		{
			name:     "Password matches input fields",
			pass:     "johndoe92",
			inputs:   []string{"John", "Doe"},
			expected: false,
		},
		{
			name:     "Secure password",
			pass:     "callgophers",
			inputs:   []string{"John", "Doe"},
			expected: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1, nil)
			got := s.Password(tt.pass, tt.inputs...)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHashAndMatch(t *testing.T) {
	cases := []struct {
		name     string
		pass     string
		expected bool
	}{
		{
			name:     "Success",
			pass:     "gamepad",
			expected: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := secure.New(1, nil)
			hash := s.Hash(tt.pass)
			assert.Equal(t, tt.expected, s.HashMatchesPassword(hash, tt.pass))
		})
	}
}

func TestToken(t *testing.T) {
	s := secure.New(1, sha1.New())
	token := "token"
	tokenized := s.Token(token)
	assert.NotEqual(t, tokenized, token)
}
