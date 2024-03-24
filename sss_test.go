package gosss

import (
	"bytes"
	"math/rand"
	"testing"
)

var examplePrivateMessage = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque nisl turpis, molestie sit amet ullamcorper sit amet, cursus in diam. Aenean urna nunc, hendrerit sed ipsum suscipit, lacinia feugiat metus. Phasellus pulvinar, tellus sit amet euismod vulputate, justo nisi finibus tellus, a ultrices odio mi vitae nibh. Duis accumsan nunc.")

func TestHideRecoverMessage(t *testing.T) {
	config := &Config{
		Shares: 4,
		Min:    3,
	}
	totalShares, err := HideMessage(examplePrivateMessage, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	candidateShares := [][]string{}
	for _, secretShares := range totalShares {
		// choose a random index to remove a share
		shares := []string{}
		index := rand.Intn(len(secretShares))
		for i, share := range secretShares {
			if i == index {
				continue
			}
			shares = append(shares, share)
		}
		candidateShares = append(candidateShares, shares)
	}
	message, err := RecoverMessage(candidateShares, config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if !bytes.Equal(message, examplePrivateMessage) {
		t.Errorf("unexpected message: %s", message)
	}
}
