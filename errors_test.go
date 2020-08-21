package fimpgo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	randomErr := errors.New("random")

	timeoutErr := errTimeout

	assert.True(t, IsTimeout(timeoutErr))
	assert.False(t, IsTimeout(randomErr))
}
