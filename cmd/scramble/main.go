package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kechako/scramble"
)

func parseNumber(s string, bits int, hex bool) (salt uint64, err error) {
	if hex {
		salt, err = strconv.ParseUint(s, 16, bits)
	} else {
		salt, err = strconv.ParseUint(s, 10, bits)
	}

	return
}

func scrambleNumber(bits int, hex bool, args []string) (int, error) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "scramble salt is not specified.")
		return 2, nil
	}

	salt, err := parseNumber(args[0], bits, hex)
	if err != nil {
		if nerr, ok := err.(*strconv.NumError); ok {
			fmt.Fprintf(os.Stderr, "scramble salt is not valid : %v\n", nerr.Err)
		} else {
			fmt.Fprintln(os.Stderr, "scramble salt is not valid.")
		}
		return 2, nil
	}
	if salt&0x01 == 0 {
		fmt.Fprintln(os.Stderr, "scramble salt is not an odd.")
		return 2, nil
	}

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "number that will be scrambled is not specified.")
		return 2, nil
	}
	num, err := parseNumber(args[1], bits, hex)
	if err != nil {
		if nerr, ok := err.(*strconv.NumError); ok {
			fmt.Fprintf(os.Stderr, "number that will be scrabled not valid : %v\n", nerr.Err)
		} else {
			fmt.Fprintln(os.Stderr, "number that will be scrabled is not valid.")
		}
		return 2, nil
	}

	var scrambled uint64
	switch bits {
	case 8:
		s := scramble.NewScrambler8WithSalt(uint8(salt))
		scrambled = uint64(s.Scramble(uint8(num)))
	case 16:
		s := scramble.NewScrambler16WithSalt(uint16(salt))
		scrambled = uint64(s.Scramble(uint16(num)))
	case 32:
		s := scramble.NewScrambler32WithSalt(uint32(salt))
		scrambled = uint64(s.Scramble(uint32(num)))
	case 64:
		s := scramble.NewScrambler64WithSalt(salt)
		scrambled = s.Scramble(num)
	}

	var format string
	if hex {
		format = fmt.Sprintf("%%0%dx\n", bits/4)
	} else {
		format = "%d\n"
	}

	fmt.Printf(format, scrambled)

	return 0, nil
}

func generateSalt(bits int, hex bool) (int, error) {
	var salt uint64
	switch bits {
	case 8:
		s, _ := scramble.GenRandomSalt8()
		salt = uint64(s)
	case 16:
		s, _ := scramble.GenRandomSalt16()
		salt = uint64(s)
	case 32:
		s, _ := scramble.GenRandomSalt32()
		salt = uint64(s)
	case 64:
		salt, _ = scramble.GenRandomSalt64()
	}

	var format string
	if hex {
		format = fmt.Sprintf("%%0%dx\n", bits/4)
	} else {
		format = "%d\n"
	}

	fmt.Printf(format, salt)

	return 0, nil
}

func main() {
	var bits int
	var genSalt bool
	var hex bool
	flag.IntVar(&bits, "bits", 32, "value bits. [8, 16, 32, 64]")
	flag.BoolVar(&genSalt, "gen", false, "generate scramble salt.")
	flag.BoolVar(&hex, "hex", false, "use hexadecimal numbers.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage : %s [options] [salt] [number]
parameters:
  salt
        scramble salt.
  number
        number that will be scrambled.

options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	switch bits {
	case 8, 16, 32, 64:
		// ok
	default:
		fmt.Fprintln(os.Stderr, "value bits is invalid.")
		os.Exit(2)
	}

	var code int
	var err error
	if genSalt {
		code, err = generateSalt(bits, hex)
	} else {
		code, err = scrambleNumber(bits, hex, flag.Args())
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}
	if code == 2 {
		flag.Usage()
	}
	if code != 0 {
		os.Exit(code)
	}
}
