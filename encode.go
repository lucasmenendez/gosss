package gosss

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// ShareSecret returns the identifier of the secret part of the share provided
// as a string. It uses the shareToStr function to encode the share. It returns
// an error if the share cannot be decoded from string.
func ShareSecret(share string) (int, error) {
	secret, _, _, err := strToShare(share)
	if err != nil {
		return 0, err
	}
	return int(secret.Int64()), nil
}

// shareToStr converts a big.Int to a string. It uses the bytes of the big.Int
func shareToStr(secret, index, share *big.Int) (string, error) {
	if secret.Cmp(big.NewInt(255)) > 0 || secret.Cmp(big.NewInt(0)) < 0 {
		return "", ErrSecretIndex
	}
	// encode the secret in a byte and append it at the end of the share
	bSecret := secret.Bytes()
	if l := len(bSecret); l > 1 {
		return "", ErrEncodeSecret
	} else if l == 0 {
		bSecret = []byte{0}
	}
	if index.Cmp(big.NewInt(255)) > 0 || index.Cmp(big.NewInt(0)) < 0 {
		return "", ErrShareIndex
	}
	bShare := share.Bytes()
	// encode the index in a byte and append it at the end of the share)
	bIndex := index.Bytes()
	if l := len(bIndex); l > 1 {
		return "", ErrEncodeIndex
	} else if l == 0 {
		bIndex = []byte{0}
	}
	fullShare := append(bShare, []byte{bIndex[0], bSecret[0]}...)
	return hex.EncodeToString(fullShare), nil
}

// strToShare converts a string to a big.Int. It uses the bytes of the string
func strToShare(s string) (*big.Int, *big.Int, *big.Int, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %w", ErrDecodeShare, err)
	}
	if len(b) < 2 {
		return nil, nil, nil, fmt.Errorf("%w: invalid lenght after decoding", ErrDecodeShare)
	}
	bSecret := b[len(b)-1]
	bIndex := b[len(b)-2]
	bShare := b[:len(b)-2]
	return new(big.Int).SetBytes([]byte{bSecret}),
		new(big.Int).SetBytes([]byte{bIndex}),
		new(big.Int).SetBytes(bShare),
		nil
}

// encodeShares function converts the x and y coordinates of the shares to
// strings. It returns the shares as strings. It uses the shareToStr function to
// encode the shares. It returns an error if the shares cannot be encoded.
func encodeShares(secrets []int, xs, ys []*big.Int) ([]string, error) {
	if len(xs) == 0 || len(ys) == 0 || len(xs) != len(ys) || len(secrets) != len(xs) {
		return nil, ErrInvalidShares
	}
	// convert the shares to strings and append them to the result
	shares := []string{}
	for i := 0; i < len(xs); i++ {
		bSecret := big.NewInt(int64(secrets[i]))
		share, err := shareToStr(bSecret, xs[i], ys[i])
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
func decodeShares(shares []string) ([]int, []*big.Int, []*big.Int, error) {
	secrets := []int{}
	xs := []*big.Int{}
	ys := []*big.Int{}
	for _, strShare := range shares {
		secret, index, share, err := strToShare(strShare)
		if err != nil {
			return nil, nil, nil, err
		}
		secrets = append(secrets, int(secret.Int64()))
		xs = append(xs, index)
		ys = append(ys, share)
	}
	return secrets, xs, ys, nil
}
