// Package scramble is a package that scramble integer numbers.
package scramble

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"unsafe"
)

type Type interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

// GenRandomSalt randomly generates a salt used by scrambler.
// The salt is an unsigned integer type.
func GenRandomSalt[T Type]() (T, error) {
	var salt T
	// salt should be grater than 1 and
	// salt must be an odd number
	for salt <= 1 || salt&0x01 == 0 {
		err := binary.Read(rand.Reader, binary.BigEndian, &salt)
		if err != nil {
			return 0, fmt.Errorf("failed to read random byte: %w", err)
		}
	}

	return salt, nil
}

// GenSaltInverse returns a inverse of the salt.
// The salt and the inverse is an unsigned integer type.
func GenSaltInverse[T Type](salt T) (T, error) {
	inv, err := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), getBits[T]())
	if err != nil {
		return 0, err
	}
	return T(inv.Uint64()), nil
}

func getBits[T Type]() int {
	var zero T
	return int(unsafe.Sizeof(zero)) * 8
}

func genSaltInverse(salt *big.Int, bits int) (*big.Int, error) {
	if salt == nil {
		return nil, errors.New("salt is nil")
	}

	switch bits {
	case 8, 16, 32, 64:
		// ok
	default:
		return nil, errors.New("invalid bits")
	}

	if salt.Bit(0) == 0 {
		return nil, errors.New("salt is not an odd number")
	}

	mod := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	return new(big.Int).ModInverse(salt, mod), nil
}

// Scrambler scrambles unsigned integers.
type Scrambler[T Type] struct {
	salt    T
	inv     T
	reverse func(T) T
}

// NewScrambler returns a new *Scramble[T] with random salt and its inverse.
func NewScrambler[T Type]() (*Scrambler[T], error) {
	salt, err := GenRandomSalt[T]()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}

	return NewScramblerWithSalt(salt)
}

// NewScramblerWithSalt returns a new *Scramble[T] with the salt and its inverse.
func NewScramblerWithSalt[T Type](salt T) (*Scrambler[T], error) {
	inv, err := GenSaltInverse(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt inverse: %w", err)
	}

	return &Scrambler[T]{
		salt:    salt,
		inv:     inv,
		reverse: getReverse[T](),
	}, nil
}

func getReverse[T Type]() func(T) T {
	var v T
	switch any(v).(type) {
	case uint8:
		return func(v T) T {
			return T(bits.Reverse8(uint8(v)))
		}
	case uint16:
		return func(v T) T {
			return T(bits.Reverse16(uint16(v)))
		}
	case uint32:
		return func(v T) T {
			return T(bits.Reverse32(uint32(v)))
		}
	default:
		return func(v T) T {
			return T(bits.Reverse64(uint64(v)))
		}
	}
}

// Scramble scrambles v.
func (s *Scrambler[T]) Scramble(v T) T {
	v *= s.salt

	v = s.reverse(v)

	v *= s.inv

	return v
}
