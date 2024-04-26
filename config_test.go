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
}

func Test_maxSecretPartSize(t *testing.T) {
	c := Config{Prime: big.NewInt(13)}
	if c.MaxMessageLen() != 0 {
		t.Errorf("expected 0, got %d", c.MaxMessageLen())
	}
	c.Prime = big.NewInt(1000000000000000000)
	if c.MaxMessageLen() != 7 {
		t.Errorf("expected 7, got %d", c.MaxMessageLen())
	}
}
