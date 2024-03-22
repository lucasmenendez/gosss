package gosss

import (
	"encoding/hex"
	"math/big"
)

// solvePolynomial solves a polynomial with coefficients coeffs for x
// and returns the result. It follows the formula:
// f(x) = a0 + a1*x + a2*x^2 + ... + an*x^n
// where a0, a1, ..., an are the coefficients.
// It uses the Horner's method to avoid the calculation of the powers of x.
// It uses the prime number defined in the package as finite field.
func solvePolynomial(coeffs []*big.Int, x, prime *big.Int) *big.Int {
	result := new(big.Int)
	result.Set(coeffs[0])
	for i := 1; i < len(coeffs); i++ {
		result.Mul(result, x)
		result.Add(result, coeffs[i])
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
func shareToStr(share *big.Int) string {
	return hex.EncodeToString(share.Bytes())
}

// strToShare converts a string to a big.Int. It uses the bytes of the string
func strToShare(s string) *big.Int {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return new(big.Int).SetBytes(b)
}
