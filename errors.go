package gosss

import "fmt"

var (
	// config
	ErrConfigShares = fmt.Errorf("wrong number of shares")
	ErrConfigMin    = fmt.Errorf("wrong minimum number of shares")
	// encode
	ErrSecretIndex   = fmt.Errorf("the index of the secret must fit in a byte (0-255)")
	ErrShareIndex    = fmt.Errorf("the index of the share must fit in a byte (0-255)")
	ErrEncodeIndex   = fmt.Errorf("error encoding index share, it must fit in a byte")
	ErrEncodeSecret  = fmt.Errorf("error encoding secret share, it must fit in a byte")
	ErrDecodeShare   = fmt.Errorf("error decoding share")
	ErrInvalidShares = fmt.Errorf("invalid shares")
	// math
	ErrReadingRandom = fmt.Errorf("error reading random number")
	// sss
	ErrRequiredConfig = fmt.Errorf("configuration is required")
	ErrEncodeMessage  = fmt.Errorf("error encoding message")
)
