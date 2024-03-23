package gosss

import (
	"math/rand"
	"testing"
)

const examplePrivateMessage = "56b919ced4ba8ee15e4f73c76120064629f056407d42c58a07106dd8f5096cd0"

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
	for _, share := range totalShares {
		t.Log(share, len(share))
	}
	// get 30 shares randomly of 36 and recover the message
	shares := []string{}
	choosen := map[string]int{}
	for len(choosen) < config.Min {
		// random number between 0 and 35
		idx := rand.Intn(config.Shares)
		_, ok := choosen[totalShares[idx]]
		for ok {
			idx = rand.Intn(config.Shares)
			_, ok = choosen[totalShares[idx]]
		}
		choosen[totalShares[idx]] = idx
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
