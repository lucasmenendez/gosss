package gosss

import "fmt"

var (
	// config
	ErrConfigShares = fmt.Errorf("wrong number of shares, it must be between 1 and the maximum number of shares")
	ErrConfigMin    = fmt.Errorf("wrong minimum number of shares, it must be between 2 and the number of shares minus 1")
	ErrConfigOp     = fmt.Errorf("unknown operation")
	// encode
	ErrShareIndex    = fmt.Errorf("the index must fit in a byte (0-255)")
	ErrDecodeShare   = fmt.Errorf("error decoding share")
	ErrInvalidShares = fmt.Errorf("invalid shares")
	// math
	ErrReadingRandom = fmt.Errorf("error reading random number")
	// sss
	ErrRequiredConfig = fmt.Errorf("configuration is required")
	ErrEncodeMessage  = fmt.Errorf("error encoding message")
)
