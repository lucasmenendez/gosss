package gosss

import (
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
}
