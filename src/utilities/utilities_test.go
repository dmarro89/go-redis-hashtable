package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	generateRandomBytes(16)

	assert.NotNil(t, randomBytes)
	assert.Len(t, randomBytes, 16)
}

func TestGetRandomBytes(t *testing.T) {
	resultBytes := GetRandomBytes()
	assert.NotEmpty(t, resultBytes)
}
