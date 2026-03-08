package security

import "errors"

var (
	ErrInvalidTokenFormat     = errors.New("invalid token format")
	ErrUnsupportedSigningAlgo = errors.New("unsupported signing algorithm")
)
