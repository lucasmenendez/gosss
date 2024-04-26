package gosss

import "fmt"

var (
	// config
	ErrRequiredConfig     = fmt.Errorf("configuration is required")
	ErrConfigShares       = fmt.Errorf("wrong number of shares")
	ErrConfigMin          = fmt.Errorf("wrong minimum number of shares")
	ErrConfigNoPrime      = fmt.Errorf("no prime provided")
	ErrConfigInvalidPrime = fmt.Errorf("invalid prime provided")
	ErrMessageTooLong     = fmt.Errorf("the message cannot be hidden with the prime provided")
	// encode
	ErrShareTooLong = fmt.Errorf("error encoding share, it is too long")
	ErrInvalidShare = fmt.Errorf("error decoding share, it is invalid")
	// math
	ErrReadingRandom = fmt.Errorf("error reading random number")
)
