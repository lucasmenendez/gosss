package gosss

import (
	"math/big"
	"testing"
)

func Test_shareToStrStrToShare(t *testing.T) {
	// generate 10 random big.Int and convert them to string
	for i := 0; i < 10; i++ {
		idx := big.NewInt(int64(i))
		rand, err := randBigInt()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		shareStr, err := shareToStr(idx, rand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		index, shareBack, err := strToShare(shareStr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if index.Cmp(idx) != 0 {
			t.Errorf("unexpected index: %d", index)
		}
		if rand.Cmp(shareBack) != 0 {
			t.Errorf("unexpected share: %s", shareStr)
		}
	}
}

func Test_encodeDecodeShares(t *testing.T) {
	xs := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
	}
	ys := []*big.Int{
		big.NewInt(100),
		big.NewInt(200),
		big.NewInt(300),
		big.NewInt(400),
	}
	encodedShares, err := encodeShares(xs, ys)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	decodedXs, decodedYs, err := decodeShares(encodedShares)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(xs) != len(decodedXs) || len(ys) != len(decodedYs) {
		t.Errorf("unexpected shares length")
		return
	}
	for i := 0; i < len(xs); i++ {
		if xs[i].Cmp(decodedXs[i]) != 0 {
			t.Errorf("unexpected x coordinate")
			return
		}
		if ys[i].Cmp(decodedYs[i]) != 0 {
			t.Errorf("unexpected y coordinate")
			return
		}
	}
}
