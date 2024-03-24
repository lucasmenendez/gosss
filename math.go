package gosss

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// randBigInt generates a random big.Int number. It uses the crypto/rand package
// to generate the random number. It returns an error if the random number
// cannot be generated.
func randBigInt() (*big.Int, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingRandom, err)
	}
	// convert the bytes to an int64 and ensure it is non-negative
	randomInt := int64(binary.BigEndian.Uint64(b[:])) & (1<<63 - 1)
	// scale down the random int to the range [0, max)
	return big.NewInt(randomInt % math.MaxInt64), nil
}

// calcCoeffs function generates the coefficients for the polynomial. It takes
// the secret and the number of coefficients to generate. It returns the
// coefficients as a list of big.Int. It returns an error if the coefficients
// cannot be generated. The secret is the first coefficient of the polynomial,
// the rest of the coefficients are random numbers.
func calcCoeffs(secret *big.Int, k int) ([]*big.Int, error) {
	// calculate k-1 random coefficients
	randCoeffs := make([]*big.Int, k-1)
	for i := 0; i < len(randCoeffs); i++ {
		randCoeff, err := randBigInt()
		if err != nil {
			return nil, err
		}
		randCoeffs[i] = randCoeff
	}
	// include secret as the first coefficient and return the coefficients
	return append([]*big.Int{secret}, randCoeffs...), nil
}

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

// calcShares function calculates the shares of the polynomial for the given
// coefficients and the number of shares to generate. It returns the x and y
// coordinates of the shares. The x coordinates are the index of the share and
// the y coordinates are the share itself. It uses the prime number to perform
// the modular operation in the finite field. It returns an error if the shares
// cannot be calculated. It skips the x = 0 coordinate because it is the secret
// itself.
func calcShares(coeffs []*big.Int, shares int, prime *big.Int) ([]*big.Int, []*big.Int) {
	// calculate shares solving the polynomial for x = {1, shares}, x = 0 is the
	// secret
	var xs, yx []*big.Int
	for i := 0; i < shares; i++ {
		x := big.NewInt(int64(i + 1))
		y := solvePolynomial(coeffs, x, prime)
		xs = append(xs, x)
		yx = append(yx, y)
	}
	return xs, yx
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
