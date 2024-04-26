package gosss

import (
	"bytes"
	"math/rand"
	"testing"
)

var examplePrivateMessage = []byte("Lorem ipsum.")

func TestHideRecoverMessage(t *testing.T) {
	config := &Config{
		Shares: 33,
		Min:    22,
	}
	totalShares, err := HideMessage(examplePrivateMessage, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	// get some shares randomly of the total and recover the message
	shares := []string{}
	chosen := map[string]int{}
	for len(chosen) < config.Min {
		// random number between 0 and 35
		idx := rand.Intn(config.Shares)
		_, ok := chosen[totalShares[idx]]
		for ok {
			idx = rand.Intn(config.Shares)
			_, ok = chosen[totalShares[idx]]
		}
		chosen[totalShares[idx]] = idx
		shares = append(shares, totalShares[idx])
	}
	message, err := RecoverMessage(shares, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if !bytes.Equal(message, examplePrivateMessage) {
		t.Errorf("unexpected message: %s", message)
	}
}
