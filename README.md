# scramble

[![GoDoc](https://godoc.org/github.com/kechako/scramble?status.svg)](https://godoc.org/github.com/kechako/scramble)

scramble is a golang package that scramble integer numbers.

## Installation

```console
go get github.com/kechako/scramble
```

## Usage

```golang
package main

import (
	"fmt"

	"github.com/kechako/scramble"
)

func main() {
	// scramble salt is randomly generated
	// use NewScramblerWithSalt if you want to specify a salt
	s, err := scramble.NewScrambler[uint32]()
	if err != nil {
		log.Fatal(err)
	}

	scrambled := s.Scramble(1234)
	fmt.Println(scrambled)
	// 4085920800

	unscrambled := s.Scramble(scrambled)
	fmt.Println(unscrambled)
	// 1234
}
```
