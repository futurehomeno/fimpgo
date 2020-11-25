package fimpgo

import (
	"errors"
)

var (
	errTimeout   = errors.New("request timed out")
	errSubscribe = errors.New("subscription failed")
	errPublish   = errors.New("publishing failed")
)

func IsTimeout(err error) bool {
	return errors.Is(err, errTimeout)
}
