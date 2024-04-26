package gosss

import (
	"math/big"
	"testing"
)

func Test_prepare(t *testing.T) {
	c := Config{}
	c.prepare()
	if c.Prime.Cmp(DefaultPrime) != 0 {
		t.Errorf("expected %v, got %v", DefaultPrime, c.Prime)
	}
	newPrime := big.NewInt(1003)
	c.Prime = new(big.Int).Set(newPrime)
	c.prepare()
	if c.Prime.Cmp(newPrime) != 0 {
		t.Errorf("expected %v, got %v", newPrime, c.Prime)
	}
	if err := c.ValidPrime(); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Prime = DefaultPrime
	if err := c.ValidPrime(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	c.Prime = nil
	if err := c.ValidPrime(); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestValidPrime(t *testing.T) {
	c := Config{}
	c.prepare()
	if err := c.ValidPrime(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	c.Prime = nil
	if err := c.ValidPrime(); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Prime = big.NewInt(1002)
	if err := c.ValidPrime(); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestValidConfig(t *testing.T) {
	c := Config{}
	c.prepare()
	if err := c.ValidConfig([]byte{}); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Shares = MinShares - 1
	c.Min = MinMinShares
	if err := c.ValidConfig([]byte{}); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Shares = MinShares
	c.Min = MinMinShares - 1
	if err := c.ValidConfig([]byte{}); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Shares = MinShares + 1
	c.Min = MinMinShares - 1
	if err := c.ValidConfig([]byte{}); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Shares = MinShares
	c.Min = MinMinShares
	c.Prime = nil
	if err := c.ValidConfig([]byte{}); err == nil {
		t.Errorf("expected error, got nil")
	}
	c.Prime = new(big.Int).SetUint64(10007)
	if err := c.ValidConfig([]byte("12345")); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMaxMessageLen(t *testing.T) {
	c := Config{Prime: big.NewInt(13)}
	if c.MaxMessageLen() != 0 {
		t.Errorf("expected 0, got %d", c.MaxMessageLen())
	}
	c.Prime = big.NewInt(1000000000000000000)
	if c.MaxMessageLen() != 7 {
		t.Errorf("expected 7, got %d", c.MaxMessageLen())
	}
}
