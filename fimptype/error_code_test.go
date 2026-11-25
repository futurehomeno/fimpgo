package fimptype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodes(t *testing.T) {
	for code, want := range map[ErrorCodeT]string{
		ErrorCodeBadRequest: "bad_request", ErrorCodeOutOfRange: "out_of_range", ErrorCodeNotSupported: "not_supported",
		ErrorCodeNotFound: "not_found", ErrorCodeBusy: "busy", ErrorCodeUnauthorized: "unauthorized",
		ErrorCodeDeviceRefused: "device_refused", ErrorCodeUnavailable: "unavailable", ErrorCodeRateLimited: "rate_limited",
		ErrorCodeInternal:      "internal",
		ErrorCodeT("eco_mode"): "eco_mode",
	} {
		assert.Equal(t, want, code.Str())
		assert.Equal(t, want, string(code))
	}
}
