package fimpgo

import (
	"errors"
)

var (
	errTimeout = errors.New("request timed out")
)

func IsTimeout(err error) bool {
	return errors.Is(err, errTimeout)
}
