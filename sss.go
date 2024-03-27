package gosss

import (
	"math/big"
	"slices"
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
	// split the message to a list of big.Int to be used as shamir secrets
	secrets := encodeMessage(message, conf.maxSecretPartSize())
	if len(secrets) == 0 {
		return nil, ErrEncodeMessage
	}
	// check if the configuration is valid for the number of secrets
	if err := conf.checkParts(len(secrets)); err != nil {
		return nil, err
	}
	// generate the shares for each secret and return them encoded as strings
	hSecrets := []int{}
	hXs := []*big.Int{}
	hYs := []*big.Int{}
	for i, secret := range secrets {
		// calculate random coefficients for the polynomial
		coeffs, err := calcCoeffs(secret, conf.minByPart())
		if err != nil {
			return nil, err
		}
		// calculate the shares with the polynomial and the prime number
		xs, yx := calcShares(coeffs, conf.sharesByPart(), conf.Prime)
		// append the shares to the result
		for range xs {
			hSecrets = append(hSecrets, i+1)
		}
		hXs = append(hXs, xs...)
		hYs = append(hYs, yx...)
	}
	// convert the shares to strings and append them to the result
	hShares, err := encodeShares(hSecrets, hXs, hYs)
	if err != nil {
		return nil, err
	}
	return hShares, nil
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
	// convert shares to big.Ints points coordinates
	secrets, xs, xy, err := decodeShares(inputs)
	if err != nil {
		return nil, err
	}
	// recover the secret parts from the shares
	parts, err := recoverParts(secrets, xs, xy, conf)
	if err != nil {
		return nil, err
	}
	// decode the message from the parts and return it
	return decodeMessage(parts), nil
}

func recoverParts(secrets []int, xs, ys []*big.Int, conf *Config) ([]*big.Int, error) {
	// check if all slices provided have the same length and non zero length
	if len(xs) != len(ys) || len(secrets) != len(xs) || len(secrets) == 0 {
		return nil, ErrInvalidShares
	}
	// recover the secret parts from the shares, group the shares by secret
	// and store the unique secrets to avoid duplicate calculations
	uniqueSecrets := []int{}
	xsBySecret := map[int][]*big.Int{}
	ysBySecret := map[int][]*big.Int{}
	for i, s := range secrets {
		_, exists := xsBySecret[s]
		if !exists {
			uniqueSecrets = append(uniqueSecrets, s)
		}
		xsBySecret[s] = append(xsBySecret[s], xs[i])
		ysBySecret[s] = append(ysBySecret[s], ys[i])
	}
	// sort the secrets to preserve the order of the message parts
	slices.Sort(uniqueSecrets)
	parts := []*big.Int{}
	for _, s := range uniqueSecrets {
		xs := xsBySecret[s]
		ys := ysBySecret[s]
		// check if the secret has the same number of coordinates
		if len(xs) != len(ys) || len(xs) == 0 {
			return nil, ErrInvalidShares
		}
		// calculate the secret part using the Lagrange interpolation, the
		// secret is the first coefficient of the polynomial (x = 0)
		part := lagrangeInterpolation(xs, ys, conf.Prime, big.NewInt(0))
		parts = append(parts, part)
	}
	return parts, nil
}
