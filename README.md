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

	"github.com/lucasmenendez/gosss"
)

func main() {
	// create a configuration with 8 shares and 7 minimum shares to recover the
	// message
	config := &gosss.Config{
		Shares: 8,
		Min:    7,
	}
	// hide a message with the defined configuration
	shares, err := gosss.HideMessage("secret", config)
	if err != nil {
		log.Fatalf("error hiding message: %v", err)
	}
	// print every share and exclude one share to test the recovery
	excluded := 3
	requiredShares := []string{}
	for i, s := range shares {
		log.Printf("share: %s", s)
		if i != excluded {
			requiredShares = append(requiredShares, s)
		}
	}
	// recover the message with the required shares
	message, err := gosss.RecoverMessage(requiredShares, nil)
	if err != nil {
		log.Fatalf("error recovering message: %v", err)
	}
	log.Printf("recovered message: %s", message)
}

```