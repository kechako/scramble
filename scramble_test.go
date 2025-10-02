package scramble

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_GenRandomSalt8(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt, err := GenRandomSalt[uint8]()
		if err != nil {
			t.Error(err)
		}
		if salt <= 1 {
			t.Errorf("GenRandomSalt[uint8]() => %d, want > 1", salt)
		}
		if salt&0x01 == 0 {
			t.Errorf("GenRandomSalt[uint8]() => %d, want odd number", salt)
		}
	}
}

func Test_GenRandomSalt16(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt, err := GenRandomSalt[uint16]()
		if err != nil {
			t.Error(err)
		}
		if salt <= 1 {
			t.Errorf("GenRandomSalt[uint16]() => %d, want > 1", salt)
		}
		if salt&0x01 == 0 {
			t.Errorf("GenRandomSalt[uint16]() => %d, want odd number", salt)
		}
	}
}

func Test_GenRandomSalt32(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt, err := GenRandomSalt[uint32]()
		if err != nil {
			t.Error(err)
		}
		if salt <= 1 {
			t.Errorf("GenRandomSalt[uint32]() => %d, want > 1", salt)
		}
		if salt&0x01 == 0 {
			t.Errorf("GenRandomSalt[uint32]() => %d, want odd number", salt)
		}
	}
}

func Test_GenRandomSalt64(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt, err := GenRandomSalt[uint64]()
		if err != nil {
			t.Error(err)
		}
		if salt <= 1 {
			t.Errorf("GenRandomSalt[uint64]() => %d, want > 1", salt)
		}
		if salt&0x01 == 0 {
			t.Errorf("GenRandomSalt[uint64]() => %d, want odd number", salt)
		}
	}
}

func Test_genSaltInverse(t *testing.T) {
	var err error

	_, err = genSaltInverse(nil, 8)
	if err == nil {
		t.Error("genSaltInverse(nil, 8) : must return error")
	}

	// bits is invalid
	_, err = genSaltInverse(big.NewInt(11), 7)
	if err == nil {
		t.Error("genSaltInverse(big.NewInt(11), 7) : must return error")
	}

	// salt is not an odd
	_, err = genSaltInverse(big.NewInt(10), 8)
	if err == nil {
		t.Error("genSaltInverse(big.NewInt(10), 8) : must return error")
	}
}

var testSalts8 = []uint8{
	0x01,
	0x13,
	0x57,
	0xa5,
	0xff,
}

func Test_GenSaltInverse8(t *testing.T) {
	for _, salt := range testSalts8 {
		name := fmt.Sprintf("salt=0x%02x", salt)
		t.Run(name, func(t *testing.T) {
			inv, err := GenSaltInverse[uint8](salt)
			if err != nil {
				t.Fatal(err)
			}
			if inv == 0 {
				t.Errorf("GenSaltInverse[uint8](0x%02x) should not return zero", salt)
			}
			if salt*inv != 1 {
				t.Errorf("salt(0x%02x) * inverse(0x%02x) != 1", salt, inv)
			}
		})
	}
}

var testSalts16 = []uint16{
	0x0001,
	0x0013,
	0x0577,
	0x1889,
	0xa471,
	0xffff,
}

func Test_GenSaltInverse16(t *testing.T) {
	for _, salt := range testSalts16 {
		name := fmt.Sprintf("salt=0x%04x", salt)
		t.Run(name, func(t *testing.T) {
			inv, err := GenSaltInverse[uint16](salt)
			if err != nil {
				t.Fatal(err)
			}
			if inv == 0 {
				t.Errorf("GenSaltInverse[uint16](0x%04x) should not return zero", salt)
			}
			if salt*inv != 1 {
				t.Errorf("salt(0x%04x) * inverse(0x%04x) != 1", salt, inv)
			}
		})
	}
}

const salt32 = 0x1ca7bc5b

var testSalts32 = []uint32{
	0x00000001,
	0x10758387,
	0x57893597,
	0x99af043b,
	0xcc886341,
	0xffffffff,
}

func Test_GenSaltInverse32(t *testing.T) {
	want := uint32(0x6b5f13d3)
	inv, err := GenSaltInverse[uint32](salt32)
	if err != nil {
		t.Fatal(err)
	}
	if inv != want {
		t.Errorf("GenSaltInverse32(0x%08x) => 0x%08x, want 0x%08x", salt32, inv, want)
	}

	for _, salt := range testSalts32 {
		name := fmt.Sprintf("salt=0x%08x", salt)
		t.Run(name, func(t *testing.T) {
			inv, err := GenSaltInverse[uint32](salt)
			if err != nil {
				t.Fatal(err)
			}
			if inv == 0 {
				t.Errorf("GenSaltInverse[uint32](0x%08x) should not return zero", salt)
			}
			if salt*inv != 1 {
				t.Errorf("salt(0x%08x) * inverse(0x%08x) != 1", salt, inv)
			}
		})
	}
}

var testSalts64 = []uint64{
	0x0000000000000001,
	0x148754987fda87cd,
	0x567a838e914c8187,
	0xa984fccd34ea6735,
	0xee99718afc32432b,
	0xffffffffffffffff,
}

func Test_GenSaltInverse64(t *testing.T) {
	for _, salt := range testSalts64 {
		name := fmt.Sprintf("salt=0x%016x", salt)
		t.Run(name, func(t *testing.T) {
			inv, err := GenSaltInverse[uint64](salt)
			if err != nil {
				t.Fatal(err)
			}
			if inv == 0 {
				t.Errorf("GenSaltInverse[uint64](0x%016x) should not return zero", salt)
			}
			if salt*inv != 1 {
				t.Errorf("salt(0x%016x) * inverse(0x%016x) != 1", salt, inv)
			}
		})
	}
}

func Test_NewScrambler8(t *testing.T) {
	s, err := NewScrambler[uint8]()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScrambler[uint8] should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv, err := GenSaltInverse[uint8](s.salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler8WithSalt(t *testing.T) {
	const salt = 101
	s, err := NewScramblerWithSalt[uint8](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScramblerWithSalt[uint8] should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv, err := GenSaltInverse[uint8](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

var scrambleTestValues8 = []uint8{
	0x01,
	0x1f,
	0x2a,
	0x4b,
	0x4d,
	0x52,
	0x66,
	0x93,
	0xa1,
	0xbc,
	0xef,
	0xff,
}

func Test_Scrambler8(t *testing.T) {
	s, err := NewScrambler[uint8]()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range scrambleTestValues8 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler16(t *testing.T) {
	s, err := NewScrambler[uint16]()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScrambler[uint16] should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv, err := GenSaltInverse[uint16](s.salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler16WithSalt(t *testing.T) {
	const salt = 101
	s, err := NewScramblerWithSalt[uint16](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScramblerWithSalt[uint16] should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv, err := GenSaltInverse[uint16](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

var scrambleTestValues16 = []uint16{
	0x0001,
	0x1439,
	0x1ba7,
	0x1dd4,
	0x3bc0,
	0x40db,
	0x5281,
	0x577b,
	0x5843,
	0x89a5,
	0xc2ad,
	0xffff,
}

func Test_Scrambler16(t *testing.T) {
	s, err := NewScrambler[uint16]()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range scrambleTestValues16 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler32(t *testing.T) {
	s, err := NewScrambler[uint32]()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScrambler[uint32] should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv, err := GenSaltInverse[uint32](s.salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler32WithSalt(t *testing.T) {
	const salt = 101
	s, err := NewScramblerWithSalt[uint32](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScramblerWithSalt[uint32] should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv, err := GenSaltInverse[uint32](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

var scrambleTestValues32 = []uint32{
	0x00000001,
	0x219aaaa7,
	0x2308a461,
	0x312ace7a,
	0x4bbc5273,
	0x4bc60edf,
	0x4ec7c0ee,
	0x55921886,
	0x6736e59e,
	0x71aee33c,
	0x76a36406,
	0xffffffff,
}

func Test_Scrambler32(t *testing.T) {
	s, err := NewScrambler[uint32]()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range scrambleTestValues32 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler64(t *testing.T) {
	s, err := NewScrambler[uint64]()
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScrambler[uint64] should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv, err := GenSaltInverse[uint64](s.salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler64WithSalt(t *testing.T) {
	const salt = 101
	s, err := NewScramblerWithSalt[uint64](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("NewScramblerWithSalt[uint64] should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv, err := GenSaltInverse[uint64](salt)
	if err != nil {
		t.Fatal(err)
	}
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

var scrambleTestValues64 = []uint64{
	0x0000000000000001,
	0x00d2a43b44b7c6f3,
	0x1426fa332a3ae4ce,
	0x2825cd66ae1d8dec,
	0x3e84f9e63c2a2719,
	0x44c61f1d0efa4c47,
	0x56d802d15e389aa1,
	0x58986084f9695020,
	0x5ea7f4aa52875750,
	0x61713a59c142fb0e,
	0x7230deafa29a104d,
	0xffffffffffffffff,
}

func Test_Scrambler64(t *testing.T) {
	s, err := NewScrambler[uint64]()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range scrambleTestValues64 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}
