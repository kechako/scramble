// Package scramble is a package that scramble integer numbers.
package scramble

import (
	"crypto/aes"
	"crypto/rand"
	"strconv"
	"unsafe"

	"gitlab.com/ubiqsecurity/ubiq-go/v2/structured"
)

// A Type is a constraint for unsigned integer types.
type Type interface {
	~uint32 | ~uint64
}

// GenerateKey generates a random key of the given size in bytes.
// The size must be either 16, 24, or 32.
func GenerateKey(size int) ([]byte, error) {
	switch size {
	default:
		return nil, aes.KeySizeError(size)
	case 16, 24, 32:
	}

	key := make([]byte, size)
	_, err := rand.Read(key)
	return key, err
}

// Scrambler scrambles unsigned integers.
type Scrambler[T Type] struct {
	ff1 *structured.FF1
}

// NewScrambler returns a new *Scramble[T].
// The key argument must be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func NewScrambler[T Type](key []byte) (*Scrambler[T], error) {
	switch k := len(key); k {
	default:
		return nil, aes.KeySizeError(k)
	case 16, 24, 32:
	}

	ff1, err := structured.NewFF1(key[:], nil, 0, 0, 16)
	if err != nil {
		return nil, err
	}

	return &Scrambler[T]{
		ff1: ff1,
	}, nil
}

// Scramble scrambles the given value v.
func (s *Scrambler[T]) Scramble(v T) (T, error) {
	var buf [16]byte
	hex := appendHex(buf[:0], v)

	e, err := s.ff1.Encrypt(string(hex), nil)
	if err != nil {
		return 0, err
	}
	vv, err := strconv.ParseUint(e, 16, bitLen[T]())
	if err != nil {
		return 0, err
	}
	return T(vv), nil
}

// Unscramble unscrambles the given value v.
func (s *Scrambler[T]) Unscramble(v T) (T, error) {
	var buf [16]byte
	hex := appendHex(buf[:0], v)

	e, err := s.ff1.Decrypt(string(hex), nil)
	if err != nil {
		return 0, err
	}
	vv, err := strconv.ParseUint(e, 16, bitLen[T]())
	if err != nil {
		return 0, err
	}
	return T(vv), nil
}

const hexTable = "0123456789abcdef"

func appendHex[T Type](dst []byte, v T) []byte {
	n := len(dst)
	l := hexLen[T]()
	dst = dst[:n+l]

	for i := l - 1; i >= 0; i-- {
		dst[n+i] = hexTable[v&0x0f]
		v >>= 4
	}

	return dst
}

func bitLen[T Type]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 8
}

func hexLen[T Type]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 2
}
