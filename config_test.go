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
	}
	if err := noMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	oneMin := &Config{
		Shares: 8,
		Min:    1,
	}
	if err := oneMin.prepare(hideOp); err == nil {
		t.Errorf("expected error")
	}
	if err := oneMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	wrongMin := &Config{
		Shares: 8,
		Min:    9,
	}
	if err := wrongMin.prepare(hideOp); err == nil {
		t.Errorf("expected error")
	}
	if err := wrongMin.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	noShares := &Config{
		Shares: 0,
		Min:    7,
	}
	if err := noShares.prepare(hideOp); err == nil {
		t.Errorf("expected error")
	}
	if err := noShares.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	maxShares := &Config{
		Shares: 257,
		Min:    7,
	}
	if err := maxShares.prepare(hideOp); err == nil {
		t.Errorf("expected error")
	}
	if err := maxShares.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	basicConf := &Config{
		Shares: 8,
		Min:    7,
	}
	if err := basicConf.prepare(hideOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if basicConf.Prime.Cmp(DefaultPrime) != 0 {
		t.Errorf("unexpected prime number: %v", basicConf.Prime)
	}
	basicConf.Prime = big.NewInt(1003)
	if err := basicConf.prepare(recoverOp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if basicConf.Prime.Cmp(big.NewInt(1003)) != 0 {
		t.Errorf("unexpected prime number: %v", basicConf.Prime)
	}
}
