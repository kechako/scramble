package scramble

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
	"testing"
)

var scrambler64Tests = []uint64{
	0xab4aa1356dfd281f,
	0x91fa7913fa867858,
	0xc50d478b0f84c119,
	0x6d7a2e13901870e7,
	0xf30fe4a6e4204e8d,
	0xdf331b0e2e9c882b,
	0xae67fe2405dd04ee,
	0xe8696a8c03cea176,
	0x873f89dfadb46c9c,
	0x7c8d817633338e15,
	0xa0938bed2fefacba,
	0xb1fdbc162eaf3e16,
	0x584e3ed38a608f57,
	0x8c03a611f662c855,
	0xe539101de31db7f5,
	0x50a4378e483aa004,
	0xd7367fc4d4aaa066,
	0x7db4c76d037bd8b7,
	0x160759765a2251f8,
	0x0643f6f2e7c092bc,
}

func TestScrambler64(t *testing.T) {
	key, err := GenerateKey(16)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewScrambler[uint64](key)
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range scrambler64Tests {
		scrambled, err := s.Scramble(input)
		if err != nil {
			t.Fatal(err)
		}
		output, err := s.Unscramble(scrambled)
		if err != nil {
			t.Fatal(err)
		}
		if output != input {
			t.Fatalf("expected %d, got %d", input, output)
		}
	}
}

var scrambler32Tests = []uint32{
	0x1370b296,
	0x73880210,
	0xe25f118d,
	0x8372882c,
	0xdffd3b3b,
	0xdf161a66,
	0xfebd91ef,
	0xc1573c5c,
	0xeb607af1,
	0xa2ed01ab,
	0x3cb23b16,
	0x672d1d89,
	0x80293ccd,
	0x56b93b5f,
	0x4afaa3d7,
	0xb34f95cc,
	0x7ebf84ef,
	0xa666e3a8,
	0xaa783753,
	0xe0515ce2,
}

func TestScrambler32(t *testing.T) {
	key, err := GenerateKey(16)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewScrambler[uint32](key)
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range scrambler32Tests {
		scrambled, err := s.Scramble(input)
		if err != nil {
			t.Fatal(err)
		}
		output, err := s.Unscramble(scrambled)
		if err != nil {
			t.Fatal(err)
		}
		if output != input {
			t.Fatalf("expected %d, got %d", input, output)
		}
	}
}

func TestAppendHex(t *testing.T) {
	var buf [16]byte
	for range 20 {
		input := rand.Uint64()
		want := fmt.Sprintf("%016x", input)
		got := string(appendHex(buf[:0], input))
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	}
}

func ExampleScrambler_Scramble() {
	key, err := hex.DecodeString("bfdaaa8ac7c7b7c84c788a3cb52a7e12")
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewScrambler[uint32](key)
	if err != nil {
		log.Fatal(err)
	}

	scrambled, err := s.Scramble(1234)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(scrambled)
	// Output: 432742918
}

func ExampleScrambler_Unscramble() {
	key, err := hex.DecodeString("bfdaaa8ac7c7b7c84c788a3cb52a7e12")
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewScrambler[uint32](key)
	if err != nil {
		log.Fatal(err)
	}

	unscrambled, err := s.Unscramble(432742918)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(unscrambled)
	// Output: 1234
}
