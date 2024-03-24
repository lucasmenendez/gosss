package gosss

import "math/big"

// 12th Mersenne Prime (2^127 - 1)
var DefaultPrime *big.Int = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(127), nil), big.NewInt(1))

// maxShares defines the maximum number of shares that can be generated
const maxShares = 256

// operation type defines the operation to perform with the configuration
type operation int

const (
	// hideOp constant defines the operation to hide a message
	hideOp operation = iota
	// recoverOp constant defines the operation to recover a message
	recoverOp
)

// Config struct defines the configuration for the Shamir Secret Sharing
// algorithm. It includes the number of shares to generate, the minimum number
// of shares to recover the secret, and the prime number to use as finite field.
type Config struct {
	Shares int
	Min    int
	Prime  *big.Int
}

// prepare checks if the configuration is valid for the operation to perform and
// sets the default prime number if it is not defined.
func (c *Config) prepare(op operation) error {
	switch op {
	case hideOp:
		// a config is valid for hide a message if the number of shares are
		// greater than 0 and lower or equal to the maximum number of shares,
		// and the minimum number of shares is greater than 1 and lower than
		// the number of shares
		if c.Shares <= 0 || c.Shares > maxShares {
			return ErrConfigShares
		}
		if c.Min <= 1 || c.Min >= c.Shares {
			return ErrConfigMin
		}
	case recoverOp:
		// for recover a message no checks are needed unless the prime number is
		// not defined so break the switch and set the default prime number if
		// it is needed
		break
	default:
		return ErrConfigOp
	}
	// if the prime number is not defined it will use the default prime number
	if c.Prime == nil {
		c.Prime = DefaultPrime
	}
	return nil
}

// maxSecretPartSize returns the maximum size of the secret part that can be
// hidden in a share, it is the size of the prime number in bytes minus 1, to
// ensure the secret part is smaller than the prime number.
func (c *Config) maxSecretPartSize() int {
	return len(c.Prime.Bytes()) - 1
}
