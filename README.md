# gosss - Go Shamir's Secret Sharing

[![GoDoc](https://godoc.org/github.com/lucasmenendez/gosss?status.svg)](https://godoc.org/github.com/lucasmenendez/gosss)
[![Build Status](https://github.com/lucasmenendez/gosss/actions/workflows/main.yml/badge.svg)](https://github.com/lucasmenendez/gosss/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasmenendez/gosss)](https://goreportcard.com/report/github.com/lucasmenendez/gosss)
[![license](https://img.shields.io/github/license/lucasmenendez/gosss)](LICENSE)


`gosss` is a Go library implementing the Shamir's Secret Sharing algorithm, a cryptographic method for splitting a secret into multiple parts. This implementation allows for secure sharing and reconstruction of secrets in a distributed system.

⚠️ This is a *for fun* implementation, it is not ready for use in a production system. ⚠️


## Getting Started

### Installation

To use `gosss` in your Go project, install it using `go get`:

```sh
go get github.com/lucasmenendez/gosss
```

### Usage
Here's a simple example of how to use gosss to split and recover a secret:

```go
package main

import (
	"log"
	"math/rand"

	"github.com/lucasmenendez/gosss"
)

func main() {
	// create a configuration with 8 shares and 7 minimum shares to recover the
	// message
	config := &gosss.Config{
		Shares: 4,
		Min:    3,
	}
	// hide a message with the defined configuration
	msg := "688641b753f1c97526d6a767058a80fd6c6519f5bdb0a08098986b0478c8502b"
	log.Printf("message to hide: %s", msg)
	totalShares, err := gosss.HideMessage([]byte(msg), config)
	if err != nil {
		log.Fatalf("error hiding message: %v", err)
	}
	// print every share and exclude one share to test the recovery
	requiredShares := [][]string{}
	for _, secretShares := range totalShares {
		log.Printf("shares: %v", secretShares)
		// choose a random share to exclude
		index := rand.Intn(len(secretShares))
		shares := []string{}
		for i, share := range secretShares {
			if i == index {
				continue
			}
			shares = append(shares, share)
		}
		requiredShares = append(requiredShares, shares)
	}
	// recover the message with the required shares
	message, err := gosss.RecoverMessage(requiredShares, nil)
	if err != nil {
		log.Fatalf("error recovering message: %v", err)
	}
	log.Printf("recovered message: %s", string(message))
}
```