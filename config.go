package gosss

import (
	"fmt"
	"math"
	"math/big"
)

// bn254 𝔽r
var DefaultPrime, _ = new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

// maxShares defines the maximum number of shares that can be generated, indexed
// using 2 bytes (255^2)
const maxShares = 65536

// ConfigLimits returns the limits for the number of shares and minimum shares
// to recover the original message based on the message size. It returns the
// minimum number of shares, the maximum number of shares, the minimum number
// of shares to recover the secret, and the maximum number of shares to recover
// the secret. The message is divided into parts based on the size of the prime
// number used as finite field, the number of parts is used to calculate the
// limits.
func ConfigLimits(message []byte) (int, int, int, int) {
	c := &Config{Prime: DefaultPrime}
	secrets := encodeMessage(message, c.maxSecretPartSize())
	nSecrets := len(secrets)
	minShares := nSecrets * 3
	minMin := nSecrets * 2
	maxMin := maxShares - nSecrets
	return minShares, maxShares, minMin, maxMin
}

// Config struct defines the configuration for the Shamir Secret Sharing
// algorithm. It includes the number of shares to generate, the minimum number
// of shares to recover the secret, and the prime number to use as finite field.
type Config struct {
	Shares int
	Min    int
	Prime  *big.Int
	nParts float64
}

// prepare sets the prime number to use as finite field if it is not defined.
func (c *Config) prepare() {
	// if the prime number is not defined it will use the default prime number
	if c.Prime == nil {
		c.Prime = DefaultPrime
	}
}

// checkParts validates the number of shares and minimum shares to recover the
// original message based on the number of parts the message is divided into. It
// returns an error if the configuration is invalid. If the configuration is
// valid, it sets the number of parts the message is divided into in the config
// struct and returns nil. The config must fit the following conditions:
//   - The number of shares must be between 3 * number of parts and 65536, a
//     constat defined by a maximum of 2 bytes to index the shares (255^2).
//   - The number of shares must be divisible by the number of parts, to ensure
//     the shares are generated correctly.
//   - The minimum number of shares to recover the secret must be between
//     2 * number of parts and number of shares - number of parts.
func (c *Config) checkParts(n int) error {
	minShares := n * 3
	if c.Shares < minShares || c.Shares > maxShares {
		return fmt.Errorf("%w (%d-%d)", ErrConfigShares, minShares, maxShares)
	}
	if c.Shares%n != 0 {
		return fmt.Errorf("%w, must be divisible by %d (number of message parts)", ErrConfigShares, n)
	}
	minMin := n * 2
	maxMin := c.Shares - n
	if c.Min < minMin || c.Min > maxMin {
		return fmt.Errorf("%w (%d-%d)", ErrConfigMin, minMin, maxMin)
	}
	c.nParts = float64(n)
	return nil
}

// maxSecretPartSize returns the maximum size of the secret part that can be
// hidden in a share, it is the size of the prime number in bytes minus 1, to
// ensure the secret part is smaller than the prime number.
func (c *Config) maxSecretPartSize() int {
	return len(c.Prime.Bytes()) - 1
}

// minByPart returns the minimum number of shares to recover the secret part, it
// is the minimum number of shares divided by the number of parts.
func (c *Config) minByPart() int {
	if c.nParts == 0 {
		return 0
	}
	return int(math.Floor(float64(c.Min) / c.nParts))
}

// sharesByPart returns the number of shares to generate for each part of the
// message, it is the number of shares divided by the number of parts.
func (c *Config) sharesByPart() int {
	if c.nParts == 0 {
		return 0
	}
	return int(math.Floor(float64(c.Shares) / c.nParts))
}
