package gosss

import (
	"math/big"
	"math/rand"
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
		share := new(big.Int).Mul(DefaultPrime, big.NewInt(rand.Int63()))
		shareStr := shareToStr(share)
		shareBack := strToShare(shareStr)
		if share.Cmp(shareBack) != 0 {
			t.Errorf("unexpected share: %s", shareStr)
		}
	}

}
