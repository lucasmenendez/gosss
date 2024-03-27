package gosss

import (
	"bytes"
	"testing"
)

var examplePrivateMessage = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque nisl turpis, molestie sit amet ullamcorper sit amet, cursus in diam. Aenean urna nunc, hendrerit sed ipsum suscipit, lacinia feugiat metus. Phasellus pulvinar, tellus sit amet euismod vulputate, justo nisi finibus tellus, a ultrices odio mi vitae nibh. Duis accumsan nunc.")

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

	completedSecrets := map[int][]string{}
	candidateShares := []string{}
	for _, share := range totalShares {
		secret, _, _, err := strToShare(share)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		iSecret := int(secret.Int64())
		curret, ok := completedSecrets[iSecret]
		if !ok || len(curret) < config.minByPart() {
			completedSecrets[iSecret] = append(completedSecrets[iSecret], share)
			candidateShares = append(candidateShares, share)
		}
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
