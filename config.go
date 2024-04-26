package gosss

import (
	"math/big"
)

// bn254 𝔽r
var DefaultPrime, _ = new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)

const (
	MinShares    = 3
	MinMinShares = MinShares - 1
)

// Config struct defines the configuration for the Shamir Secret Sharing
// algorithm. It includes the number of shares to generate, the minimum number
// of shares to recover the secret, and the prime number to use as finite field.
type Config struct {
	Shares int
	Min    int
	Prime  *big.Int
}

// prepare sets the prime number to use as finite field if it is not defined or
// if it is smaller than 2 bytes. It uses the default prime number defined in
// the package.
func (c *Config) prepare() {
	// if the prime number is not defined or is smaller than 2 bytes, it will
	// use the default prime number
	if c.Prime == nil || len(c.Prime.Bytes()) < 2 {
		c.Prime = DefaultPrime
	}
}

// MaxMessageLen returns the maximum size of the secret that can be hidden
// in a share, it is the size of the prime number in bytes minus 1, to ensure
// the secret is smaller than the prime number.
func (c *Config) MaxMessageLen() int {
	if max := len(c.Prime.Bytes()) - 1; max > 0 {
		return max
	}
	return 0
}

// ValidPrime checks if the configuration has a valid prime number. It returns
// an error if the prime number is not defined or if it is not a prime number.
func (c *Config) ValidPrime() error {
	// check if the prime number is a prime number
	if c.Prime == nil {
		return ErrConfigNoPrime
	}
	if !c.Prime.ProbablyPrime(0) {
		return ErrConfigInvalidPrime
	}
	if len(c.Prime.Bytes()) < 2 {
		return ErrConfigInvalidPrime
	}
	return nil
}

// ValidConfig checks if the configuration is valid for the secret provided.
// It checks if the number of shares is greater than the minimum number of
// shares, if the minimum number of shares is greater than the number of shares
// less one or if it is smaller than the minimum number of shares less one, if
// the config has a valid prime number, and if the message can be hidden with
// the prime number.
func (c *Config) ValidConfig(secret []byte) error {
	// check if the number of shares is greater than the minimum number of shares
	if c.Shares < MinShares {
		return ErrConfigShares
	}
	// check if the minimum number of shares is greater than the number of shares
	// less one or if it is smaller than the minimum number of shares less one
	if c.Min > c.Shares-1 || c.Min < MinMinShares {
		return ErrConfigMin
	}
	// check if the config has a valid prime number
	if err := c.ValidPrime(); err != nil {
		return err
	}
	// check if the message can be hidden with the prime number
	if len(secret) > c.MaxMessageLen() {
		return ErrMessageTooLong
	}
	return nil
}
