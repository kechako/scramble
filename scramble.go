package scramble

import (
	"math"
	"math/big"
	"math/bits"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

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

func GenSaltInverse8(salt uint8) uint8 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 8)

	return uint8(inv.Uint64())
}

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

func GenSaltInverse16(salt uint16) uint16 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 16)

	return uint16(inv.Uint64())
}

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

func GenSaltInverse32(salt uint32) uint32 {
	inv := genSaltInverse(new(big.Int).SetUint64(uint64(salt)), 32)

	return uint32(inv.Uint64())
}

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

type Scrambler8 struct {
	salt uint8
	inv  uint8
}

func NewScrambler8() *Scrambler8 {
	salt, inv := GenRandomSalt8()

	return &Scrambler8{
		salt: salt,
		inv:  inv,
	}
}

func NewScrambler8WithSalt(salt uint8) *Scrambler8 {
	inv := GenSaltInverse8(salt)

	return &Scrambler8{
		salt: salt,
		inv:  inv,
	}
}

func (s *Scrambler8) Scramble(v uint8) uint8 {
	v *= s.salt

	v = bits.Reverse8(v)

	v *= s.inv

	return v
}

type Scrambler16 struct {
	salt uint16
	inv  uint16
}

func NewScrambler16() *Scrambler16 {
	salt, inv := GenRandomSalt16()

	return &Scrambler16{
		salt: salt,
		inv:  inv,
	}
}

func NewScrambler16WithSalt(salt uint16) *Scrambler16 {
	inv := GenSaltInverse16(salt)

	return &Scrambler16{
		salt: salt,
		inv:  inv,
	}
}

func (s *Scrambler16) Scramble(v uint16) uint16 {
	v *= s.salt

	v = bits.Reverse16(v)

	v *= s.inv

	return v
}

type Scrambler32 struct {
	salt uint32
	inv  uint32
}

func NewScrambler32() *Scrambler32 {
	salt, inv := GenRandomSalt32()

	return &Scrambler32{
		salt: salt,
		inv:  inv,
	}
}

func NewScrambler32WithSalt(salt uint32) *Scrambler32 {
	inv := GenSaltInverse32(salt)

	return &Scrambler32{
		salt: salt,
		inv:  inv,
	}
}

func (s *Scrambler32) Scramble(v uint32) uint32 {
	v *= s.salt

	v = bits.Reverse32(v)

	v *= s.inv

	return v
}

type Scrambler64 struct {
	salt uint64
	inv  uint64
}

func NewScrambler64() *Scrambler64 {
	salt, inv := GenRandomSalt64()

	return &Scrambler64{
		salt: salt,
		inv:  inv,
	}
}

func NewScrambler64WithSalt(salt uint64) *Scrambler64 {
	inv := GenSaltInverse64(salt)

	return &Scrambler64{
		salt: salt,
		inv:  inv,
	}
}

func (s *Scrambler64) Scramble(v uint64) uint64 {
	v *= s.salt

	v = bits.Reverse64(v)

	v *= s.inv

	return v
}
