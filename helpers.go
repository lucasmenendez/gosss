package gosss

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
)

// solvePolynomial solves a polynomial with coefficients coeffs for x
// and returns the result. It follows the formula:
// f(x) = a0 + a1*x + a2*x^2 + ... + an*x^n
// where a0, a1, ..., an are the coefficients.
// It uses the Horner's method to avoid the calculation of the powers of x.
// It uses the prime number defined in the package as finite field.
func solvePolynomial(coeffs []*big.Int, x, prime *big.Int) *big.Int {
	accum := big.NewInt(0)
	for i := len(coeffs) - 1; i >= 0; i-- {
		accum.Mul(accum, x)
		accum.Add(accum, coeffs[i])
		accum.Mod(accum, prime)
	}
	return accum
}

// lagrangeInterpolation calculates the Lagrange interpolation for the given
// points and a specific x value. The formula for Lagrange interpolation over a
// finite field defined by a prime is:
//
//	L(x) = sum(y_j * product((x - x_m) / (x_j - x_m) for m != j), for all j)
//
// where the operations are performed modulo a prime number to ensure results
// remain within the finite field.
func lagrangeInterpolation(xCoords, yCoords []*big.Int, prime, _x *big.Int) *big.Int {
	// initialize result
	result := big.NewInt(0)
	// temporary variables used inside the loop to avoid reallocating them, it
	// will be overwritten in each iteration with some values that are not used
	// in the next iteration
	temp1 := big.NewInt(0)
	temp2 := big.NewInt(0)
	// numerator and denominator for each point to be calculated in the loop
	numerator := big.NewInt(1)
	denominator := big.NewInt(1)
	diff := big.NewInt(0)
	// iterate over all the points
	for i := range xCoords {
		// reset the numerator and denominator for each point
		numerator.SetInt64(1)
		denominator.SetInt64(1)
		for j := range xCoords {
			// skip if i == j
			if i == j {
				continue
			}
			// calculate numerator: (x - x_m)
			temp1.Sub(_x, xCoords[j])
			numerator.Mul(numerator, temp1)
			// modular operation
			numerator.Mod(numerator, prime)
			// calculate denominator: (x_j - x_m)
			diff.Sub(xCoords[i], xCoords[j])
			// modular inverse
			temp2.ModInverse(diff, prime)
			denominator.Mul(denominator, temp2)
			// modular operation
			denominator.Mod(denominator, prime)
		}
		// combine the fraction with y_i and add it to the result
		temp1.Mul(yCoords[i], numerator) // y_i * numerator
		temp1.Mod(temp1, prime)
		temp2.Mul(temp1, denominator) // y_i * numerator / denominator
		temp2.Mod(temp2, prime)
		// add the result of the fraction to the final result
		result.Add(result, temp2)
		// ensure the result is within the field
		result.Mod(result, prime)
	}
	return result
}

// msgToBigInt converts a string to a big.Int. It uses the bytes of the string
// to create the big.Int.
func msgToBigInt(s string) *big.Int {
	return new(big.Int).SetBytes([]byte(s))
}

// bigIntToMsg converts a big.Int to a string. It uses the bytes of the big.Int
// to create the string.
func bigIntToMsg(i *big.Int) string {
	return string(i.Bytes())
}

// shareToStr converts a big.Int to a string. It uses the bytes of the big.Int
func shareToStr(index, share *big.Int) (string, error) {
	if index.Cmp(big.NewInt(255)) > 0 || index.Cmp(big.NewInt(0)) < 0 {
		return "", fmt.Errorf("the index must fit in a byte (0-255)")
	}
	bShare := share.Bytes()
	// encode the index in a byte and append it at the end of the share
	bIndex := index.Bytes()
	if len(bIndex) == 0 {
		bIndex = []byte{0}
	}
	fullShare := append(bShare, bIndex[0])
	return hex.EncodeToString(fullShare), nil
}

// strToShare converts a string to a big.Int. It uses the bytes of the string
func strToShare(s string) (*big.Int, *big.Int, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding share: %v", err)
	}
	bIndex := b[len(b)-1]
	bShare := b[:len(b)-1]
	return new(big.Int).SetBytes([]byte{bIndex}), new(big.Int).SetBytes(bShare), nil
}

// randBigInt generates a random big.Int number. It uses the crypto/rand package
// to generate the random number. It returns an error if the random number
// cannot be generated.
func randBigInt() (*big.Int, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, err
	}
	// convert the bytes to an int64 and ensure it is non-negative
	randomInt := int64(binary.BigEndian.Uint64(b[:])) & (1<<63 - 1)
	// scale down the random int to the range [0, max)
	return big.NewInt(randomInt % math.MaxInt64), nil
}
