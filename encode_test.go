package gosss

import (
	"math/big"
	"testing"
)

func Test_shareToStrStrToShare(t *testing.T) {
	// generate 10 random big.Int and convert them to string
	for i := 0; i < 10; i++ {
		x, err := randBigInt(8, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		y, err := randBigInt(8, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		shareStr, err := shareToStr(x, y)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		nx, ny, err := strToShare(shareStr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if x.Cmp(nx) != 0 {
			t.Errorf("unexpected secret: %d", x)
		}
		if y.Cmp(ny) != 0 {
			t.Errorf("unexpected share: %s", y)
		}
	}
	// test coords
	invalidX, _ := new(big.Int).SetString("1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 10)
	invalidY, _ := new(big.Int).SetString("1", 10)
	_, err := shareToStr(invalidX, invalidY)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	// test invalid share
	_, _, err = strToShare("3")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	_, _, err = strToShare("10")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
