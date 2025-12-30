package gifdecode

import "errors"

var (
	ErrTooLarge    = errors.New("gif exceeds limits")
	ErrNoFrames    = errors.New("gif has no frames")
	ErrInvalidSize = errors.New("gif has invalid size")
)
