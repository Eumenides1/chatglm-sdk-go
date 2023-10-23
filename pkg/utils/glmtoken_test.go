package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGlmToken(t *testing.T) {
	glmToken := NewGlmToken(&APICredentials{
		ApiKey:    "{id}",
		ApiSecret: "{secret}",
	})
	fmt.Printf(glmToken.Token)
	cacheToken := findTokenInCache("{id}")
	assert.Equal(t, cacheToken, glmToken.Token)
}
