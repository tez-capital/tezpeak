package constants

import "errors"

var (
	ErrInvalidListenAddress    = errors.New("invalid listen address")
	ErrInvalidWorkingDirectory = errors.New("invalid working directory")
	ErrInvalidBlockWindow      = errors.New("invalid block window")
	ErrInvalidConfigVersion    = errors.New("invalid configuration version")
	ErrInvalidConfig           = errors.New("invalid configuration")
	ErrInvalidSignerUrl        = errors.New("invalid signer url")

	ErrFailedToSignOperation      = errors.New("failed to sign operation")
	ErrFailedToCompleteOperation  = errors.New("failed to complete operation")
	ErrFailedToGetPublicKey       = errors.New("failed to get public key")
	ErrFailedToCreateRemoteSigner = errors.New("failed to create remote signer")
	ErrFailedToBroadcastOperation = errors.New("failed to broadcast operation")
	ErrDelegateNotRegistered      = errors.New("delegate not registered")
)
