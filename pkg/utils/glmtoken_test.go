package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGlmToken(t *testing.T) {
	glmToken := NewGlmToken(&APICredentials{
		ApiKey:    "",
		ApiSecret: "",
	})
	fmt.Printf(glmToken.Token)
	cacheToken := findTokenInCache("")
	assert.Equal(t, cacheToken, glmToken.Token)
}
