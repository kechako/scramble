# scramble

[![GoDoc](https://godoc.org/github.com/kechako/scramble/v2?status.svg)](https://godoc.org/github.com/kechako/scramble/v2)

scramble is a Go library that performs format-preserving scrambling and unscrambling of numeric values using the FF1 FPE algorithm.
It provides simple APIs that take and return uint32 or uint64, enabling reversible obfuscation while preserving the original numeric format.
This makes it well suited for anonymization, tokenization, and other use cases where fixed-size numeric identifiers must remain valid.

## Installation

```console
go get github.com/kechako/scramble/v2
```

## Usage

```golang
package main

import (
	"fmt"
	"log"

	scramble "github.com/kechako/scramble/v2"
)

func main() {
	// Generate a random key
	key, err := scramble.GenerateKey(16)
	if err != nil {
		log.Fatal(err)
	}

	// Create a scrambler for uint32 using the generated key
	s, err := scramble.NewScrambler[uint32](key)
	if err != nil {
		log.Fatal(err)
	}

	scrambled, err := s.Scramble(1234)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(scrambled)
	// e.g. 4085920800

	unscrambled, err := s.Unscramble(scrambled)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(unscrambled)
	// 1234
}
```
