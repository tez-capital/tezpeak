package constants

import "errors"

var (
	ErrInvalidWorkingDirectory = errors.New("invalid working directory")
	ErrInvalidBlockWindow      = errors.New("invalid block window")
	ErrInvalidConfigVersion    = errors.New("invalid configuration version")
	ErrInvalidConfig           = errors.New("invalid configuration")
)
