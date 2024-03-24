package gosss

import (
	"bytes"
	"math/big"
	"testing"
)

func Test_encodeDecodeMessage(t *testing.T) {
	var maxPartSize = len(DefaultPrime.Bytes()) - 1
	var inputMessage = []byte("688641b753f1c97526d6a767058a80fd6c6519f5bdb0a08098986b0478c8502b")
	var expectedParts = []*big.Int{
		new(big.Int).SetBytes([]byte("688641b753f1c97")),
		new(big.Int).SetBytes([]byte("526d6a767058a80")),
		new(big.Int).SetBytes([]byte("fd6c6519f5bdb0a")),
		new(big.Int).SetBytes([]byte("08098986b0478c8")),
		new(big.Int).SetBytes([]byte("502b")),
	}
	parts := encodeMessage(inputMessage, maxPartSize)
	if len(parts) != len(expectedParts) {
		t.Errorf("Expected %d parts but got %d", len(expectedParts), len(parts))
		return
	}
	for i, part := range parts {
		if part.Cmp(expectedParts[i]) != 0 {
			t.Errorf("Expected part %d to be %d but got %d", i, expectedParts[i], part)
			return
		}
	}
	if decodedMessage := decodeMessage(parts); !bytes.Equal(inputMessage, decodedMessage) {
		t.Errorf("Expected %s but got %s", inputMessage, decodedMessage)
	}
}
