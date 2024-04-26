package gosss

import (
	"fmt"
	"math/big"
)

// HideMessage generates the shares of the message using the Shamir Secret
// Sharing algorithm. It returns the shares as strings. The message is encoded
// as a big.Int and the shares are calculated solving a polynomial with random
// coefficients. The first coefficient is the encoded message. It uses the
// configuration provided in the Config struct, if the prime number is not
// defined it uses the 12th Mersenne Prime (2^127 - 1) as default. It returns
// an error if the message cannot be encoded.
func HideMessage(message []byte, conf *Config) ([]string, error) {
	// the hide operation needs the minimum number of shares and the total
	// number of shares, so if the configuration is not provided, return an
	// error
	if conf == nil {
		return nil, ErrRequiredConfig
	}
	conf.prepare()
	// validate the configuration for the message provided
	if err := conf.ValidConfig(message); err != nil {
		return nil, err
	}
	// calculate k random coefficients for the polynomial, where k is the
	// minimum number of shares less one (the secret is the first coefficient)
	coeffs, err := calcCoeffs(new(big.Int).SetBytes(message), conf.Prime, conf.Min)
	if err != nil {
		return nil, err
	}
	// calculate the shares with the polynomial and the prime number
	xs, ys := calcShares(coeffs, conf.Shares, conf.Prime)
	// encode the shares
	shares := []string{}
	for i := 0; i < len(xs); i++ {
		share, err := shareToStr(xs[i], ys[i])
		if err != nil {
			return nil, fmt.Errorf("error encoding shares: %w", err)
		}
		shares = append(shares, share)
	}
	return shares, nil
}

// RecoverMessage recovers the message from the shares using the Shamir Secret
// Sharing algorithm. It returns the message as a string. The shares are given
// as strings. It uses the configuration provided in the Config struct, if the
// prime number is not defined it uses the 12th Mersenne Prime (2^127 - 1) as
// default. It returns an error if the message cannot be recovered. The shares
// include the index of the share and the share itself, so the order of the
// provided shares does not matter. It decodes the points of the polynomial from
// the shares and calculates the Lagrange interpolation to recover the secret.
func RecoverMessage(inputs []string, conf *Config) ([]byte, error) {
	// the recover operation does not need the minimum number of shares or the
	// total number of shares, so if the configuration is not provided, create a
	// empty configuration before prepare the it.
	if conf == nil {
		conf = &Config{}
	}
	// prepare the configuration to recover the message
	conf.prepare()
	if err := conf.ValidPrime(); err != nil {
		return nil, err
	}
	// convert shares to big.Ints points coordinates
	xs, ys := []*big.Int{}, []*big.Int{}
	for _, input := range inputs {
		x, y, err := strToShare(input)
		if err != nil {
			return nil, fmt.Errorf("error decoding shares: %w", err)
		}
		xs = append(xs, x)
		ys = append(ys, y)
	}
	// calculate the secret using the Lagrange interpolation, the secret is the
	// first coefficient of the polynomial (x = 0)
	secret := lagrangeInterpolation(xs, ys, conf.Prime, big.NewInt(0))
	// decode the message from the secret y coord
	return secret.Bytes(), nil
}
