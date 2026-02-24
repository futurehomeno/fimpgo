package utils

import "errors"

var (
	ErrNoAvailableInterface  = errors.New("no available interface")
	ErrNoAvailableAddress    = errors.New("no available address")
	ErrNotSupportedInterface = errors.New("not supported interface")
	ErrTimeout               = errors.New("timeout")
	ErrConnectionLost        = errors.New("connection_lost")
	ErrEmptyCache            = errors.New("empty_cache")
)
