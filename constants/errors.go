package constants

import "errors"

var (
	ErrInvalidListenAddress    = errors.New("invalid listen address")
	ErrInvalidWorkingDirectory = errors.New("invalid working directory")
	ErrInvalidBlockWindow      = errors.New("invalid block window")
	ErrInvalidConfigVersion    = errors.New("invalid configuration version")
	ErrInvalidConfig           = errors.New("invalid configuration")
)
