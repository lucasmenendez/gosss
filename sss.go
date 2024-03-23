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
func HideMessage(message string, conf *Config) ([]string, error) {
	// the hide operation needs the minimum number of shares and the total
	// number of shares, so if the configuration is not provided, return an
	// error
	if conf == nil {
		return nil, fmt.Errorf("configuration is required")
	}
	// prepare the configuration to hide the message
	if err := conf.prepare(hideOp); err != nil {
		return nil, err
	}
	// encode message to big.Int
	secret := msgToBigInt(message)
	if secret == nil {
		return nil, fmt.Errorf("error encoding message")
	}
	// calculate k-1 random coefficients (k = min)
	randCoeffs := make([]*big.Int, conf.Min-1)
	for i := 0; i < len(randCoeffs); i++ {
		randCoeff, err := randBigInt()
		if err != nil {
			return nil, err
		}
		randCoeffs[i] = randCoeff
	}
	// include secret as the first coefficient
	coeffs := append([]*big.Int{secret}, randCoeffs...)
	// calculate shares solving the polynomial for x = {1, shares}, x = 0 is the
	// secret
	totalShares := make([]string, conf.Shares)
	for i := 0; i < conf.Shares; i++ {
		x := big.NewInt(int64(i + 1))
		y := solvePolynomial(coeffs, x, conf.Prime)
		share, err := shareToStr(x, y)
		if err != nil {
			return nil, err
		}
		totalShares[i] = share
	}
	return totalShares, nil
}

// RecoverMessage recovers the message from the shares using the Shamir Secret
// Sharing algorithm. It returns the message as a string. The shares are given
// as strings. It uses the configuration provided in the Config struct, if the
// prime number is not defined it uses the 12th Mersenne Prime (2^127 - 1) as
// default. It returns an error if the message cannot be recovered. The shares
// include the index of the share and the share itself, so the order of the
// provided shares does not matter. It decodes the points of the polynomial from
// the shares and calculates the Lagrange interpolation to recover the secret.
func RecoverMessage(shares []string, conf *Config) (string, error) {
	// the recover operation does not need the minimum number of shares or the
	// total number of shares, so if the configuration is not provided, create a
	// empty configuration before prepare the it.
	if conf == nil {
		conf = &Config{}
	}
	// prepare the configuration to recover the message
	if err := conf.prepare(recoverOp); err != nil {
		return "", err
	}
	// convert shares to big.Ints points coordinates
	x := make([]*big.Int, len(shares))
	y := make([]*big.Int, len(shares))
	for i, strShare := range shares {
		index, share, err := strToShare(strShare)
		if err != nil {
			return "", err
		}
		x[i] = index
		y[i] = share
	}
	result := lagrangeInterpolation(x, y, conf.Prime, big.NewInt(0))
	return bigIntToMsg(result), nil
}
