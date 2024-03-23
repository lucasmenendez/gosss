package gosss

import (
	"fmt"
	"math/big"
	"math/rand"
)

// HideMessage generates the shares of the message using the Shamir Secret
// Sharing algorithm. It returns the shares as strings. The message is encoded
// as a big.Int and the shares are calculated solving a polynomial with random
// coefficients. The first coefficient is the encoded message. It uses the
// configuration provided in the Config struct, if the prime number is not
// defined it uses the 12th Mersenne Prime (2^127 - 1) as default. It returns
// an error if the message cannot be encoded.
func HideMessage(message string, conf *Config) ([]string, error) {
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
		randCoeffs[i] = new(big.Int).Mul(conf.Prime, big.NewInt(rand.Int63()))
	}
	// include secret as the first coefficient
	coeffs := append([]*big.Int{secret}, randCoeffs...)
	// calculate shares solving the polynomial for x = {1, shares}, x = 0 is the
	// secret
	totalShares := make([]string, conf.Shares)
	for i := 0; i < conf.Shares; i++ {
		x := big.NewInt(int64(i + 1))
		y := solvePolynomial(coeffs, x, conf.Prime)
		totalShares[i] = shareToStr(y)
	}
	return totalShares, nil
}

func RecoverMessage(shares []string, conf *Config) (string, error) {
	// prepare the configuration to recover the message
	if err := conf.prepare(recoverOp); err != nil {
		return "", err
	}
	// convert shares to big.Int
	bShares := make([]*big.Int, len(shares))
	for i, share := range shares {
		bShares[i] = strToShare(share)
		if bShares[i] == nil {
			return "", fmt.Errorf("error decoding share")
		}
	}
	var result *big.Int
	// TODO: Implement the recovery of the secret using Lagrange interpolation
	return bigIntToMsg(result), nil
}
