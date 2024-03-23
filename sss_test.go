package gosss

import (
	"math/rand"
	"testing"
)

const examplePrivateMessage = "aaa"

func TestHideRecoverMessage(t *testing.T) {
	config := &Config{
		Shares: 8,
		Min:    7,
	}
	totalShares, err := HideMessage(examplePrivateMessage, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(totalShares) != config.Shares {
		t.Errorf("unexpected number of shares: %d", len(totalShares))
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
	if examplePrivateMessage != message {
		t.Errorf("unexpected message: %s", message)
	}
}
