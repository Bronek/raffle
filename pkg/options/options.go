package options

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Input          string
	Seed           int64
	SeedFromSystem bool
	Multiplier     int
	N              int
}

func (config Config) String() string {
	systemEntropy := ""
	if config.SeedFromSystem {
		systemEntropy = " (from system entropy)"
	}
	return fmt.Sprintf(" --input='%s'\n --seed=%v%s\n --multiplier=%v\n --N=%v\n",
		config.Input, config.Seed, systemEntropy, config.Multiplier, config.N)
}

func Parse(progname string, args []string) (config *Config, help string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var result Config
	flags.StringVar(
		&result.Input,
		"input",
		"",
		"Input file, expected CSV with two columns: name and number of tickets\n"+
			"Number of tickets may contain decimal part (see `multiplier` option).\n"+
			"If input file is not provided, program will use standard input.")
	flags.Int64Var(
		&result.Seed,
		"seed",
		0,
		"Random seed for shuffle operation. If not provided or set to 0,\n"+
			"program will use system entropy")
	flags.IntVar(
		&result.Multiplier,
		"multiplier",
		1,
		"Ticket size multiplier, allowed values between 1 and 10000.\n"+
			"Use 100 to handle ticket sizes with two decimal places")
	flags.IntVar(
		&result.N,
		"N",
		1,
		"Number of winning spots to produce on program output.")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	if result.Multiplier < 1 || result.Multiplier > 10_000 {
		return nil, "", fmt.Errorf("Multiplier out of range 1 to 10000: %v", result.Multiplier)
	}

	result.SeedFromSystem = false
	// Loop here because Seed=0 has special meaning and we do
	// not want it if system entropy gives us this value
	for result.Seed == 0 {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to access system entropy: %w", err)
		}
		result.Seed = int64(binary.BigEndian.Uint64(b))
		result.SeedFromSystem = true
	}
	return &result, "", nil
}

func InputFile(config *Config) (file *os.File, err error) {
	if len(config.Input) == 0 {
		return os.Stdin, nil
	}

	fd, err := os.Open(config.Input)
	if err != nil {
		return nil, fmt.Errorf("Unable to open input file: %w", err)
	}
	return fd, nil
}
