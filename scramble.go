// Package scramble is a package that scramble integer numbers.
package scramble

import (
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenRandomSalt8 randomly generates a salt and its inverse used by scrambler.
// The salt and the inverse is a 8 bits unsigned integer.
func GenRandomSalt8() (uint8, uint8) {
	var salt uint8
	for salt <= 1 {
		// salt should be grater than 1
		salt = uint8(random.Intn(math.MaxUint8) + 1)

		// salt must be an odd number
		if salt&0x01 == 0 {
			salt++
		}
	}

	return salt, GenSaltInverse8(salt)
}

// GenSaltInverse8 returns a inverse of the salt.
// The salt and the inverse is a 8 bits unsigned integer.
func GenSaltInverse8(salt uint8) uint8 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 8)

	return uint8(inv.Uint64())
}

// GenRandomSalt16 randomly generates a salt and its inverse used by scrambler.
// The salt is a 16 bits unsigned integer.
func GenRandomSalt16() (uint16, uint16) {
	var salt uint16
	for salt <= 1 {
		// salt should be grater than 1
		salt = uint16(random.Intn(math.MaxUint16) + 1)

		// salt must be an odd number
		if salt&0x01 == 0 {
			salt++
		}
	}

	return salt, GenSaltInverse16(salt)
}

// GenSaltInverse16 returns a inverse of the salt.
// The salt and the inverse is a 16 bits unsigned integer.
func GenSaltInverse16(salt uint16) uint16 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 16)

	return uint16(inv.Uint64())
}

// GenRandomSalt32 randomly generates a salt and its inverse used by scrambler.
// The salt is a 32 bits unsigned integer.
func GenRandomSalt32() (uint32, uint32) {
	var salt uint32
	for salt <= 1 {
		// salt should be grater than 1
		salt = uint32(random.Int31() + 1)

		// salt must be an odd number
		if salt&0x01 == 0 {
			salt++
		}
	}

	return salt, GenSaltInverse32(salt)
}

// GenSaltInverse32 returns a inverse of the salt.
// The salt and the inverse is a 32 bits unsigned integer.
func GenSaltInverse32(salt uint32) uint32 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 32)

	return uint32(inv.Uint64())
}

// GenRandomSalt64 randomly generates a salt and its inverse used by scrambler.
// The salt is a 64 bits unsigned integer.
func GenRandomSalt64() (uint64, uint64) {
	var salt uint64
	for salt <= 1 {
		// salt should be grater than 1
		salt = uint64(random.Int63() + 1)

		// salt must be an odd number
		if salt&0x01 == 0 {
			salt++
		}
	}

	return salt, GenSaltInverse64(salt)
}

// GenSaltInverse64 returns a inverse of the salt.
// The salt and the inverse is a 64 bits unsigned integer.
func GenSaltInverse64(salt uint64) uint64 {
	inv := genSaltInverse(new(big.Int).SetUint64(salt), 64)

	return inv.Uint64()
}

func genSaltInverse(salt *big.Int, bits int) *big.Int {
	if salt == nil {
		panic("salt is nil")
	}

	switch bits {
	case 8, 16, 32, 64:
		// ok
	default:
		panic("invalid bits")
	}

	if salt.Bit(0) == 0 {
		panic("salt is not an odd number")
	}

	mod := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	return new(big.Int).ModInverse(salt, mod)
}

// Scrambler8 scrambles 8 bit unsigned integers.
type Scrambler8 struct {
	salt uint8
	inv  uint8
}

// NewScrambler8 returns a new *Scramble8 with random salt and its inverse.
func NewScrambler8() *Scrambler8 {
	salt, inv := GenRandomSalt8()

	return &Scrambler8{
		salt: salt,
		inv:  inv,
	}
}

// NewScrambler8WithSalt returns a new *Scramble8 with the salt and its inverse.
func NewScrambler8WithSalt(salt uint8) *Scrambler8 {
	inv := GenSaltInverse8(salt)

	return &Scrambler8{
		salt: salt,
		inv:  inv,
	}
}

// Scramble scrambles v.
func (s *Scrambler8) Scramble(v uint8) uint8 {
	v *= s.salt

	v = bits.Reverse8(v)

	v *= s.inv

	return v
}

// Scrambler16 scrambles 16 bit unsigned integers.
type Scrambler16 struct {
	salt uint16
	inv  uint16
}

// NewScrambler16 returns a new *Scramble16 with random salt and its inverse.
func NewScrambler16() *Scrambler16 {
	salt, inv := GenRandomSalt16()

	return &Scrambler16{
		salt: salt,
		inv:  inv,
	}
}

// NewScrambler16WithSalt returns a new *Scramble16 with the salt and its inverse.
func NewScrambler16WithSalt(salt uint16) *Scrambler16 {
	inv := GenSaltInverse16(salt)

	return &Scrambler16{
		salt: salt,
		inv:  inv,
	}
}

// Scramble scrambles v.
func (s *Scrambler16) Scramble(v uint16) uint16 {
	v *= s.salt

	v = bits.Reverse16(v)

	v *= s.inv

	return v
}

// Scrambler32 scrambles 32 bit unsigned integers.
type Scrambler32 struct {
	salt uint32
	inv  uint32
}

// NewScrambler32 returns a new *Scramble32 with random salt and its inverse.
func NewScrambler32() *Scrambler32 {
	salt, inv := GenRandomSalt32()

	return &Scrambler32{
		salt: salt,
		inv:  inv,
	}
}

// NewScrambler32WithSalt returns a new *Scramble32 with the salt and its inverse.
func NewScrambler32WithSalt(salt uint32) *Scrambler32 {
	inv := GenSaltInverse32(salt)

	return &Scrambler32{
		salt: salt,
		inv:  inv,
	}
}

// Scramble scrambles v.
func (s *Scrambler32) Scramble(v uint32) uint32 {
	v *= s.salt

	v = bits.Reverse32(v)

	v *= s.inv

	return v
}

// Scrambler64 scrambles 64 bit unsigned integers.
type Scrambler64 struct {
	salt uint64
	inv  uint64
}

// NewScrambler64 returns a new *Scramble64 with random salt and its inverse.
func NewScrambler64() *Scrambler64 {
	salt, inv := GenRandomSalt64()

	return &Scrambler64{
		salt: salt,
		inv:  inv,
	}
}

// NewScrambler64WithSalt returns a new *Scramble64 with the salt and its inverse.
func NewScrambler64WithSalt(salt uint64) *Scrambler64 {
	inv := GenSaltInverse64(salt)

	return &Scrambler64{
		salt: salt,
		inv:  inv,
	}
}

// Scramble scrambles v.
func (s *Scrambler64) Scramble(v uint64) uint64 {
	v *= s.salt

	v = bits.Reverse64(v)

	v *= s.inv

	return v
}
