package gosss

import (
	"encoding/hex"
	"math/big"
)

// shareToStr converts a big.Int to a string. It uses the bytes of the big.Int
func shareToStr(x, y *big.Int) (string, error) {
	bx, by := x.Bytes(), y.Bytes()
	lx, ly := len(bx), len(by)
	if lx > 255 || ly > 255 {
		return "", ErrShareTooLong
	}
	fullShare := append(bx, by...)
	fullShare = append(fullShare, []byte{byte(lx), byte(ly)}...)
	return hex.EncodeToString(fullShare), nil
}

// strToShare converts a string to a big.Int. It uses the bytes of the string
func strToShare(s string) (*big.Int, *big.Int, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, ErrInvalidShare
	}
	lb := len(b)
	if len(b) < 2 {
		return nil, nil, ErrInvalidShare
	}
	lx, ly := int(b[lb-2]), int(b[lb-1])
	bx, by := b[:lx], b[lx:lx+ly]
	return new(big.Int).SetBytes(bx), new(big.Int).SetBytes(by), nil
}
