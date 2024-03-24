package gosss

import "math/big"

// HideMessage generates the shares of the message using the Shamir Secret
// Sharing algorithm. It returns the shares as strings. The message is encoded
// as a big.Int and the shares are calculated solving a polynomial with random
// coefficients. The first coefficient is the encoded message. It uses the
// configuration provided in the Config struct, if the prime number is not
// defined it uses the 12th Mersenne Prime (2^127 - 1) as default. It returns
// an error if the message cannot be encoded.
func HideMessage(message []byte, conf *Config) ([][]string, error) {
	// the hide operation needs the minimum number of shares and the total
	// number of shares, so if the configuration is not provided, return an
	// error
	if conf == nil {
		return nil, ErrRequiredConfig
	}
	// prepare the configuration to hide the message
	if err := conf.prepare(hideOp); err != nil {
		return nil, err
	}
	// split the message to a list of big.Int to be used as shamir secrets
	secrets := encodeMessage(message, conf.maxSecretPartSize())
	if len(secrets) == 0 {
		return nil, ErrEncodeMessage
	}
	// generate the shares for each secret and return them encoded as strings
	shares := [][]string{}
	for _, secret := range secrets {
		// calculate random coefficients for the polynomial
		coeffs, err := calcCoeffs(secret, conf.Min)
		if err != nil {
			return nil, err
		}
		// calculate the shares with the polynomial and the prime number
		xs, yx := calcShares(coeffs, conf.Shares, conf.Prime)
		// convert the shares to strings and append them to the result
		secretShares, err := encodeShares(xs, yx)
		if err != nil {
			return nil, err
		}
		shares = append(shares, secretShares)
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
func RecoverMessage(shares [][]string, conf *Config) ([]byte, error) {
	// the recover operation does not need the minimum number of shares or the
	// total number of shares, so if the configuration is not provided, create a
	// empty configuration before prepare the it.
	if conf == nil {
		conf = &Config{}
	}
	// prepare the configuration to recover the message
	if err := conf.prepare(recoverOp); err != nil {
		return nil, err
	}
	parts := []*big.Int{}
	for _, secretShares := range shares {
		// convert shares to big.Ints points coordinates
		xs, ys, err := decodeShares(secretShares)
		if err != nil {
			return nil, err
		}
		// calculate the secret part using the Lagrange interpolation, the
		// secret part is the y coordinate for x = 0
		result := lagrangeInterpolation(xs, ys, conf.Prime, big.NewInt(0))
		// append the secret part to the result
		parts = append(parts, result)
	}
	// decode the message from the parts and return it
	return decodeMessage(parts), nil
}
