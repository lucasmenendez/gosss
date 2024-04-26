package gosss

import (
	"math/big"
	"testing"
)

func Test_randBigInt(t *testing.T) {
	generatedRands := make(map[int64]bool)
	for i := 0; i < 100000; i++ {
		rand, err := randBigInt(8, nil)
		if err != nil {
			t.Fatalf("error generating random number: %v", err)
			return
		}
		if _, ok := generatedRands[rand.Int64()]; ok {
			t.Fatalf("duplicated random number")
			return
		}
		generatedRands[rand.Int64()] = true
	}
}

func Test_calcCoeffs(t *testing.T) {
	secret := big.NewInt(123456789)
	coeffs, err := calcCoeffs(secret, DefaultPrime, 5)
	if err != nil {
		t.Fatalf("error calculating coefficients: %v", err)
		return
	}
	if len(coeffs) != 5 {
		t.Fatalf("invalid number of coefficients")
		return
	}
	if coeffs[0].Cmp(secret) != 0 {
		t.Fatalf("invalid secret coefficient")
		return
	}
	checkedCoeffs := make(map[int64]bool)
	for i := 1; i < len(coeffs); i++ {
		if _, ok := checkedCoeffs[coeffs[i].Int64()]; ok {
			t.Fatalf("duplicated coefficient")
			return
		}
		checkedCoeffs[coeffs[i].Int64()] = true
	}
}

func Test_solvePolynomial(t *testing.T) {
	// f(x) = 1 + 2x + 3x^2 + 4x^3
	basicCoeffs := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
	}
	// x = 2, prime = 5
	basicX := big.NewInt(2)
	basicPrime := big.NewInt(5)
	// f(2) = 1 + 4 + 12 + 32 = 49 % 5 = 4
	basicExpected := big.NewInt(4)
	basicResult := solvePolynomial(basicCoeffs, basicX, basicPrime)
	if basicResult.Cmp(basicExpected) != 0 {
		t.Errorf("Simple polynomial failed, expected %v, got %v", basicExpected, basicResult)
	}

	// f(x) = 0
	zeroCoeffs := []*big.Int{
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
	}
	// x = 5, prime = 7
	zeroX := big.NewInt(5)
	zeroPrime := big.NewInt(7)
	// f(5) = 0
	zeroExpected := big.NewInt(0)
	zeroResult := solvePolynomial(zeroCoeffs, zeroX, zeroPrime)
	if zeroResult.Cmp(zeroExpected) != 0 {
		t.Errorf("Zero polynomial failed, expected %v, got %v", zeroExpected, zeroResult)
	}

	// f(x) = -1 - 2x - 3x^2 - 4x^3
	negativeCoeffs := []*big.Int{
		big.NewInt(-1),
		big.NewInt(-2),
		big.NewInt(-3),
		big.NewInt(-4),
	}
	// x = -2, prime = 5
	negativeX := big.NewInt(-2)
	negativePrime := big.NewInt(5)
	// f(-2) = -1 + 4 - 12 + 32 = 23 % 5 = 3
	negativeExpected := big.NewInt(3)
	negativeResult := solvePolynomial(negativeCoeffs, negativeX, negativePrime)
	if negativeResult.Cmp(negativeExpected) != 0 {
		t.Errorf("Negative polynomial failed, expected %v, got %v", negativeExpected, negativeResult)
	}

	// f(x) = 1 + 2x + 3x^2 + 4x^3 + 5x^4 + 6x^5 + 7x^6 + 8x^7
	highDegreeCoeffs := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
		big.NewInt(5),
		big.NewInt(6),
		big.NewInt(7),
		big.NewInt(8),
	}
	highDegreeX := 2
	highDegreePrime := 13
	// f(2) = 1 + 4 + 12 + 32 + 80 + 192 + 448 + 1024 = 1793 % 13 = 12
	highDegreeExpected := big.NewInt(12)
	highDegreeResult := solvePolynomial(highDegreeCoeffs, big.NewInt(int64(highDegreeX)), big.NewInt(int64(highDegreePrime)))
	if highDegreeResult.Cmp(highDegreeExpected) != 0 {
		t.Errorf("High degree polynomial failed, expected %v, got %v", highDegreeExpected, highDegreeResult)
	}
}

func Test_calcShares(t *testing.T) {
	prime := big.NewInt(5)
	// f(x) = a + bx + cx^2 + dx^3 = 5 + x + 2x^2 + 3x^3
	coeffs := []*big.Int{
		big.NewInt(5),
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
	}
	var solve = func(x, prime *big.Int) *big.Int {
		a := new(big.Int).Set(coeffs[0])
		b := new(big.Int).Mul(coeffs[1], x)
		c := new(big.Int).Mul(coeffs[2], new(big.Int).Mul(x, x))
		d := new(big.Int).Mul(coeffs[3], new(big.Int).Mul(x, new(big.Int).Mul(x, x)))
		res := new(big.Int).Add(a, new(big.Int).Add(b, new(big.Int).Add(c, d)))
		return new(big.Int).Mod(res, prime)
	}

	k := len(coeffs)
	xs, ys := calcShares(coeffs, k, prime)
	if len(xs) != k {
		t.Fatalf("invalid number of x coordinates")
		return
	}
	if len(ys) != k {
		t.Fatalf("invalid number of y coordinates")
		return
	}
	for i, x := range xs {
		y := solve(x, prime)
		if y.Cmp(ys[i]) != 0 {
			t.Errorf("invalid y coordinate, expected %v, got %v", y, ys[i])
		}
	}
}

func Test_lagrangeInterpolation(t *testing.T) {
	prime := big.NewInt(17)
	// f(x) = (6 + x + 2x^2 + 3x^3) % 17
	coeffs := []*big.Int{
		big.NewInt(6),
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
	}
	xs, ys := calcShares(coeffs, len(coeffs), prime)
	// x = 0; f(0) = 6 % 17 = 6
	x0 := big.NewInt(0)
	y0 := big.NewInt(6)
	result0 := lagrangeInterpolation(xs, ys, prime, x0)
	if result0.Cmp(y0) != 0 {
		t.Errorf("x = 0 failed, expected %v, got %v", y0, result0)
		return
	}

	// x = 3; f(3) = 6 + 3 + 18 + 81 = 108 % 17 = 6
	x3 := big.NewInt(3)
	y3 := big.NewInt(6)
	result3 := lagrangeInterpolation(xs, ys, prime, x3)
	if result3.Cmp(y3) != 0 {
		t.Errorf("x = 3 failed, expected %v, got %v", y3, result3)
		return
	}
	// x = 4; f(4) = 6 + 4 + 32 + 192 = 234 % 17 = 13
	x4 := big.NewInt(4)
	y4 := big.NewInt(13)
	result4 := lagrangeInterpolation(xs, ys, prime, x4)
	if result4.Cmp(y4) != 0 {
		t.Errorf("x = 4 failed, expected %v, got %v", y4, result4)
	}
}
