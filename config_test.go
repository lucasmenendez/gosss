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
	newPrime := big.NewInt(103)
	c.Prime = new(big.Int).Set(newPrime)
	c.prepare()
	if c.Prime.Cmp(newPrime) != 0 {
		t.Errorf("expected %v, got %v", newPrime, c.Prime)
	}
}

func Test_checkParts(t *testing.T) {
	// valid configuration
	validConfig := Config{Shares: 12, Min: 8, Prime: big.NewInt(13)}
	if err := validConfig.checkParts(4); err != nil {
		t.Errorf("Valid Configuration: Expected no error, got %v", err)
		return
	}
	// shares less than minimum
	lessThanMinSharesConfig := Config{Shares: 2, Min: 1, Prime: big.NewInt(13)}
	if err := lessThanMinSharesConfig.checkParts(1); err == nil {
		t.Errorf("shares less than minimum: Expected error, got none")
		return
	}
	// shares greater than maximum
	greaterThanMaxSharesConfig := Config{Shares: maxShares + 1, Min: 20000, Prime: big.NewInt(13)}
	if err := greaterThanMaxSharesConfig.checkParts(10000); err == nil {
		t.Errorf("shares Greater than maximum: Expected error, got none")
		return
	}
	// hares Not Divisible by parts
	sharesNotDivisibleConfig := Config{Shares: 10, Min: 5, Prime: big.NewInt(13)}
	if err := sharesNotDivisibleConfig.checkParts(3); err == nil {
		t.Errorf("shares Not Divisible by parts: Expected error, got none")
		return
	}
	// inimum shares less than allowed
	minSharesLessThanAllowedConfig := Config{Shares: 12, Min: 3, Prime: big.NewInt(13)}
	if err := minSharesLessThanAllowedConfig.checkParts(4); err == nil {
		t.Errorf("Minimum shares less than allowed: Expected error, got none")
		return
	}
	// inimum shares Greater than allowed
	minSharesGreaterThanAllowedConfig := Config{Shares: 12, Min: 10, Prime: big.NewInt(13)}
	if err := minSharesGreaterThanAllowedConfig.checkParts(4); err == nil {
		t.Errorf("Minimum shares Greater than allowed: Expected error, got none")
		return
	}
}

func Test_maxSecretPartSize(t *testing.T) {
	c := Config{Prime: big.NewInt(13)}
	if c.maxSecretPartSize() != 0 {
		t.Errorf("expected 0, got %d", c.maxSecretPartSize())
	}
	c.Prime = big.NewInt(1000000000000000000)
	if c.maxSecretPartSize() != 7 {
		t.Errorf("expected 7, got %d", c.maxSecretPartSize())
	}
}

func Test_minByPart(t *testing.T) {
	c := Config{Min: 8}
	if c.minByPart() != 0 {
		t.Errorf("expected 0, got %d", c.minByPart())
	}
	c.nParts = 4
	if c.minByPart() != 2 {
		t.Errorf("expected 2, got %d", c.minByPart())
	}
	c.nParts = 0
	if c.minByPart() != 0 {
		t.Errorf("expected 0, got %d", c.minByPart())
	}
}

func Test_sharesByPart(t *testing.T) {
	c := Config{Shares: 12}
	if c.sharesByPart() != 0 {
		t.Errorf("expected 0, got %d", c.sharesByPart())
	}
	c.nParts = 4
	if c.sharesByPart() != 3 {
		t.Errorf("expected 3, got %d", c.sharesByPart())
	}
	c.nParts = 0
	if c.sharesByPart() != 0 {
		t.Errorf("expected 0, got %d", c.sharesByPart())
	}
}