package scramble

import (
	"math/big"
	"testing"
)

var testSalts8 = []uint8{
	0x01,
	0x13,
	0x57,
	0xa5,
	0xff,
}

func Test_genSaltInverse(t *testing.T) {
	// test panic
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("genSaltInverse(nil, 8) : panic must happen")
			}
		}()
		// salt is nil
		genSaltInverse(nil, 8)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("genSaltInverse(big.NewInt(11), 7) : panic must happen")
			}
		}()
		// bits is invalid
		genSaltInverse(big.NewInt(11), 7)
	}()

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("genSaltInverse(big.NewInt(10), 8) : panic must happen")
			}
		}()
		// salt is not an odd
		genSaltInverse(big.NewInt(10), 8)
	}()
}

func Test_GenSaltInverse8(t *testing.T) {
	for _, salt := range testSalts8 {
		inv := GenSaltInverse8(salt)
		if inv == 0 {
			t.Errorf("GenSaltInverse8(0x%02x) should not return zero", salt)
		}
		if salt*inv != 1 {
			t.Errorf("salt(0x%02x) * inverse(0x%02x) != 1", salt, inv)
		}
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
		inv := GenSaltInverse16(salt)
		if inv == 0 {
			t.Errorf("GenSaltInverse16(0x%04x) should not return zero", salt)
		}
		if salt*inv != 1 {
			t.Errorf("salt(0x%04x) * inverse(0x%04x) != 1", salt, inv)
		}
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
	inv := GenSaltInverse32(salt32)
	if inv != want {
		t.Errorf("GenSaltInverse32(0x%08x) => 0x%08x, want 0x%08x", salt32, inv, want)
	}

	for _, salt := range testSalts32 {
		inv := GenSaltInverse32(salt)
		if inv == 0 {
			t.Errorf("GenSaltInverse32(0x%08x) should not return zero", salt)
		}
		if salt*inv != 1 {
			t.Errorf("salt(0x%08x) * inverse(0x%08x) != 1", salt, inv)
		}
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
		inv := GenSaltInverse64(salt)
		if inv == 0 {
			t.Errorf("GenSaltInverse64(0x%016x) should not return zero", salt)
		}
		if salt*inv != 1 {
			t.Errorf("salt(0x%016x) * inverse(0x%016x) != 1", salt, inv)
		}
	}
}

func Test_NewScrambler8(t *testing.T) {
	s := NewScrambler8()
	if s == nil {
		t.Fatal("NewScrambler8 should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv := GenSaltInverse8(s.salt)
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler8WithSalt(t *testing.T) {
	const salt = 101
	s := NewScrambler8WithSalt(salt)
	if s == nil {
		t.Fatal("NewScrambler8WithSalt should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv := GenSaltInverse8(salt)
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
	s := NewScrambler8()
	for _, value := range scrambleTestValues8 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler16(t *testing.T) {
	s := NewScrambler16()
	if s == nil {
		t.Fatal("NewScrambler16 should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv := GenSaltInverse16(s.salt)
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler16WithSalt(t *testing.T) {
	const salt = 101
	s := NewScrambler16WithSalt(salt)
	if s == nil {
		t.Fatal("NewScrambler16WithSalt should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv := GenSaltInverse16(salt)
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
	s := NewScrambler16()
	for _, value := range scrambleTestValues16 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler32(t *testing.T) {
	s := NewScrambler32()
	if s == nil {
		t.Fatal("NewScrambler32 should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv := GenSaltInverse32(s.salt)
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler32WithSalt(t *testing.T) {
	const salt = 101
	s := NewScrambler32WithSalt(salt)
	if s == nil {
		t.Fatal("NewScrambler32WithSalt should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv := GenSaltInverse32(salt)
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
	s := NewScrambler32()
	for _, value := range scrambleTestValues32 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}

func Test_NewScrambler64(t *testing.T) {
	s := NewScrambler64()
	if s == nil {
		t.Fatal("NewScrambler64 should not return nil")
	}

	if s.salt == 0 {
		t.Errorf("salt should not be zero")
	}
	if s.inv == 0 {
		t.Errorf("inv should not be zero")
	}
	inv := GenSaltInverse64(s.salt)
	if s.inv != inv {
		t.Errorf("inv should be %d, but %d", inv, s.inv)
	}
}

func Test_NewScrambler64WithSalt(t *testing.T) {
	const salt = 101
	s := NewScrambler64WithSalt(salt)
	if s == nil {
		t.Fatal("NewScrambler64WithSalt should not return nil")
	}

	if s.salt != salt {
		t.Errorf("salt should be %d, but %d", salt, s.salt)
	}
	inv := GenSaltInverse64(salt)
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
	s := NewScrambler64()
	for _, value := range scrambleTestValues64 {
		scrambled := s.Scramble(value)
		unscrambled := s.Scramble(scrambled)

		if unscrambled != value {
			t.Errorf("Scramble(%d) => %d, Scramble(%d) => %d, want %d", value, scrambled, scrambled, unscrambled, value)
		}
	}
}
