package main

import (
	"fmt"
	"log"

	scramble "github.com/kechako/scramble/v2"
)

func main() {
	key, err := scramble.GenerateKey(16)
	if err != nil {
		log.Fatal(err)
	}

	s, err := scramble.NewScrambler[uint32](key)
	if err != nil {
		log.Fatal(err)
	}

	scrambled, err := s.Scramble(1234)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(scrambled)
	// 4085920800

	unscrambled, err := s.Unscramble(scrambled)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(unscrambled)
	// 1234
}
