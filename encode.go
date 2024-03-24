package gosss

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// shareToStr converts a big.Int to a string. It uses the bytes of the big.Int
func shareToStr(index, share *big.Int) (string, error) {
	if index.Cmp(big.NewInt(255)) > 0 || index.Cmp(big.NewInt(0)) < 0 {
		return "", ErrShareIndex
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
		return nil, nil, fmt.Errorf("%w: %w", ErrDecodeShare, err)
	}
	bIndex := b[len(b)-1]
	bShare := b[:len(b)-1]
	return new(big.Int).SetBytes([]byte{bIndex}), new(big.Int).SetBytes(bShare), nil
}

// encodeShares function converts the x and y coordinates of the shares to
// strings. It returns the shares as strings. It uses the shareToStr function to
// encode the shares. It returns an error if the shares cannot be encoded.
func encodeShares(xs, ys []*big.Int) ([]string, error) {
	if len(xs) == 0 || len(ys) == 0 || len(xs) != len(ys) {
		return nil, ErrInvalidShares
	}
	// convert the shares to strings and append them to the result
	shares := []string{}
	for i := 0; i < len(xs); i++ {
		share, err := shareToStr(xs[i], ys[i])
		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}
	return shares, nil
}

// decodeShares function converts the strings of the shares to x and y
// coordinates of the shares. It uses the strToShare function to decode the
// shares. It returns an error if the shares cannot be decoded.
func decodeShares(shares []string) ([]*big.Int, []*big.Int, error) {
	xs := []*big.Int{}
	ys := []*big.Int{}
	for _, strShare := range shares {
		index, share, err := strToShare(strShare)
		if err != nil {
			return nil, nil, err
		}
		xs = append(xs, index)
		ys = append(ys, share)
	}
	return xs, ys, nil
}
