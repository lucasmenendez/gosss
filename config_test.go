package gosss

import (
	"math/big"
	"testing"
)

func Test_prepare(t *testing.T) {
	noMin := &Config{
		Shares: 8,
		Min:    0,
	}
	if err := noMin.prepare(hideOp); err == nil {
		t.Errorf("expected error")
		return
	}
	if err := noMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	oneMin := &Config{
		Shares: 8,
		Min:    1,
	}
	if err := oneMin.prepare(hideOp); err == nil {
		t.Errorf("expected error")
		return
	}
	if err := oneMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	wrongMin := &Config{
		Shares: 8,
		Min:    9,
	}
	if err := wrongMin.prepare(hideOp); err == nil {
		t.Errorf("expected error")
		return
	}
	if err := wrongMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	noShares := &Config{
		Shares: 0,
		Min:    7,
	}
	if err := noShares.prepare(hideOp); err == nil {
		t.Errorf("expected error")
		return
	}
	if err := noShares.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	maxShares := &Config{
		Shares: 257,
		Min:    7,
	}
	if err := maxShares.prepare(hideOp); err == nil {
		t.Errorf("expected error")
		return
	}
	if err := maxShares.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	basicConf := &Config{
		Shares: 8,
		Min:    7,
	}
	if err := basicConf.prepare(hideOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if basicConf.Prime.Cmp(DefaultPrime) != 0 {
		t.Errorf("unexpected prime number: %v", basicConf.Prime)
		return
	}
	basicConf.Prime = big.NewInt(1003)
	if err := basicConf.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if basicConf.Prime.Cmp(big.NewInt(1003)) != 0 {
		t.Errorf("unexpected prime number: %v", basicConf.Prime)
		return
	}
}
