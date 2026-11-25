package fimptype

type ErrorCodeT string

const (
	ErrorCodeBadRequest    ErrorCodeT = "bad_request"
	ErrorCodeOutOfRange    ErrorCodeT = "out_of_range"
	ErrorCodeNotSupported  ErrorCodeT = "not_supported"
	ErrorCodeNotFound      ErrorCodeT = "not_found"
	ErrorCodeBusy          ErrorCodeT = "busy"
	ErrorCodeUnauthorized  ErrorCodeT = "unauthorized"
	ErrorCodeDeviceRefused ErrorCodeT = "device_refused"
	ErrorCodeUnavailable   ErrorCodeT = "unavailable"
	ErrorCodeRateLimited   ErrorCodeT = "rate_limited"
	ErrorCodeInternal      ErrorCodeT = "internal"
)

func (c ErrorCodeT) Str() string {
	return string(c)
}
