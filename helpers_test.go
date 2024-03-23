package gosss

import (
	"math/big"
	"testing"
)

func Test_encodeDecodeMsg(t *testing.T) {
	encodedPrivateMsg := msgToBigInt(examplePrivateMessage)
	if encodedPrivateMsg == nil {
		t.Errorf("unexpected nil encoded string")
		return
	}
	decodedPrivateMsg := bigIntToMsg(encodedPrivateMsg)
	if examplePrivateMessage != decodedPrivateMsg {
		t.Errorf("unexpected decoded string: %s", decodedPrivateMsg)
	}
}

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
