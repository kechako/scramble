package main

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	scramble "github.com/kechako/scramble/v2"
	cli "github.com/urfave/cli/v3"
)

func scrambleCommand(ctx context.Context, c *cli.Command) error {
	keyStr := c.StringArg("key")
	if keyStr == "" {
		return cli.Exit(errors.New("scramble random key is not specified"), 2)
	}
	keyEnc := c.String("keyenc")
	key, err := decodeKey(keyStr, keyEnc)
	if err != nil {
		return cli.Exit(fmt.Errorf("scramble random key is not valid: %v", err), 2)
	}

	numStr := c.StringArg("number")
	encode := c.String("encode")
	bits := c.Int("bits")
	num, err := decodeNumber(numStr, bits, encode)
	if err != nil {
		return cli.Exit(fmt.Errorf("number is not valid: %v", err), 2)
	}

	var scrambled uint64
	switch bits {
	case 32:
		scrambled, err = scrambleNumber(key, uint32(num))
	case 64:
		scrambled, err = scrambleNumber(key, num)
	default:
		return cli.Exit(errors.New("bit size must be either 32 or 64"), 2)
	}
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to scramble number: %v", err), 1)
	}

	result, err := encodeNumber(scrambled, bits, encode)
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to encode scrambled number: %v", err), 1)
	}

	fmt.Println(result)

	return nil
}

func encodeNumber(n uint64, bits int, encode string) (string, error) {
	if encode == "none" {
		return strconv.FormatUint(n, 10), nil
	}

	l := bits / 8
	var buf [8]byte
	b := buf[:l]
	if bits == 32 {
		binary.BigEndian.PutUint32(b, uint32(n))
	} else {
		binary.BigEndian.PutUint64(b, n)
	}

	switch encode {
	case "hex":
		return hex.EncodeToString(b), nil
	case "base32":
		return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
	case "base64":
		return base64.RawURLEncoding.EncodeToString(b), nil
	default:
		return "", fmt.Errorf("unsupported encode format: %s", encode)
	}
}

func decodeNumber(s string, bits int, encode string) (uint64, error) {
	if encode == "none" {
		return strconv.ParseUint(s, 10, bits)
	}

	var b []byte
	var err error
	switch encode {
	case "hex":
		b, err = hex.DecodeString(s)
	case "base32":
		b, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
	case "base64":
		b, err = base64.RawURLEncoding.DecodeString(s)
	default:
		return 0, fmt.Errorf("unsupported encode format: %s", encode)
	}
	if err != nil {
		return 0, err
	}

	if bits == 32 {
		return uint64(binary.BigEndian.Uint32(b)), nil
	} else {
		return binary.BigEndian.Uint64(b), nil
	}
}

func scrambleNumber[T scramble.Type](key []byte, v T) (uint64, error) {
	s, err := scramble.NewScrambler[T](key)
	if err != nil {
		return 0, err
	}
	scrambled, err := s.Scramble(v)
	if err != nil {
		return 0, err
	}
	return uint64(scrambled), nil
}

func unscrambleCommand(ctx context.Context, c *cli.Command) error {
	keyStr := c.StringArg("key")
	if keyStr == "" {
		return cli.Exit(errors.New("unscramble random key is not specified"), 2)
	}
	keyEnc := c.String("keyenc")
	key, err := decodeKey(keyStr, keyEnc)
	if err != nil {
		return cli.Exit(fmt.Errorf("unscramble random key is not valid: %v", err), 2)
	}

	numStr := c.StringArg("number")
	encode := c.String("encode")
	bits := c.Int("bits")
	num, err := decodeNumber(numStr, bits, encode)
	if err != nil {
		return cli.Exit(fmt.Errorf("number is not valid: %v", err), 2)
	}

	var scrambled uint64
	switch bits {
	case 32:
		scrambled, err = unscrambleNumber(key, uint32(num))
	case 64:
		scrambled, err = unscrambleNumber(key, num)
	default:
		return cli.Exit(errors.New("bit size must be either 32 or 64"), 2)
	}
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to unscramble number: %v", err), 1)
	}

	result, err := encodeNumber(scrambled, bits, encode)
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to encode scrambled number: %v", err), 1)
	}

	fmt.Println(result)

	return nil
}

func unscrambleNumber[T scramble.Type](key []byte, v T) (uint64, error) {
	s, err := scramble.NewScrambler[T](key)
	if err != nil {
		return 0, err
	}
	unscrambled, err := s.Unscramble(v)
	if err != nil {
		return 0, err
	}
	return uint64(unscrambled), nil
}

func genCommand(ctx context.Context, c *cli.Command) error {
	keySize := c.Int("size")
	encode := c.String("encode")
	switch encode {
	case "hex", "base32", "base64":
		// ok
	default:
		return cli.Exit(fmt.Errorf("unsupported encode format: %s", encode), 2)
	}

	key, err := scramble.GenerateKey(keySize)
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to generate random key: %v", err), 1)
	}

	encoded, err := encodeKey(key, encode)
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to encode generated key: %v", err), 1)
	}

	fmt.Println(encoded)

	return nil
}

func encodeKey(key []byte, encode string) (string, error) {
	switch encode {
	case "hex":
		return hex.EncodeToString(key), nil
	case "base32":
		return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(key), nil
	case "base64":
		return base64.RawURLEncoding.EncodeToString(key), nil
	default:
		return "", fmt.Errorf("unsupported encode format: %s", encode)
	}
}

func decodeKey(s string, encode string) ([]byte, error) {
	switch encode {
	case "hex":
		return hex.DecodeString(s)
	case "base32":
		return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
	case "base64":
		return base64.RawURLEncoding.DecodeString(s)
	default:
		return nil, fmt.Errorf("unsupported encode format: %s", encode)
	}
}

func main() {
	cmd := &cli.Command{
		Name:  "scramble",
		Usage: "scramble random numbers using format-preserving encryption.",
		Commands: []*cli.Command{
			{
				Name:  "gen",
				Usage: "generate random key.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "size",
						Aliases: []string{"s"},
						Usage:   "key size in bytes. the size must be either 16, 24, or 32.",
						Value:   16,
					},
					&cli.StringFlag{
						Name:    "encode",
						Aliases: []string{"e"},
						Usage:   "how to encode the generated key. hex: hexadecimal (default), base32: base32, base64: base64.",
						Value:   "hex",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Action: genCommand,
			},
			{
				Name:  "scramble",
				Usage: "scramble a number.",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "key",
						UsageText: "key",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
					&cli.StringArg{
						Name:      "number",
						UsageText: "number",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bits",
						Aliases: []string{"b"},
						Usage:   "Bit size of the number. 32 or 64.",
						Value:   32,
					},
					&cli.StringFlag{
						Name:    "encode",
						Aliases: []string{"e"},
						Usage:   "how to encode the number. none: not encode (default), hex: hexadecimal, base32: base32, base64: base64.",
						Value:   "none",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
					&cli.StringFlag{
						Name:    "keyenc",
						Aliases: []string{"k"},
						Usage:   "how to encode the generated key. hex: hexadecimal (default), base32: base32, base64: base64.",
						Value:   "hex",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Action: scrambleCommand,
			},
			{
				Name:  "unscramble",
				Usage: "unscramble a number.",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "key",
						UsageText: "key",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
					&cli.StringArg{
						Name:      "number",
						UsageText: "number",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bits",
						Aliases: []string{"b"},
						Usage:   "Bit size of the number. 32 or 64.",
						Value:   32,
					},
					&cli.StringFlag{
						Name:    "encode",
						Aliases: []string{"e"},
						Usage:   "how to encode the number. none: not encode (default), hex: hexadecimal, base32: base32, base64: base64.",
						Value:   "none",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
					&cli.StringFlag{
						Name:    "keyenc",
						Aliases: []string{"k"},
						Usage:   "how to encode the generated key. hex: hexadecimal (default), base32: base32, base64: base64.",
						Value:   "hex",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Action: unscrambleCommand,
			},
		},
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := cmd.Run(ctx, os.Args); err != nil {
		code := 1

		var coder cli.ExitCoder
		if errors.As(err, &coder) {
			code = coder.ExitCode()
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		if code != 0 {
			os.Exit(code)
		}
	}
}
