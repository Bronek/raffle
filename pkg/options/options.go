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
	Seed           int64
	SeedFromSystem bool
	Multiplier     int
	N              int
	args           []string
}

func (config Config) String() string {
	return fmt.Sprintf("seed=%v system=%t multiplier=%v N=%v args=%+v",
		config.Seed, config.SeedFromSystem, config.Multiplier, config.N, config.args)
}

func Parse(progname string, args []string) (config *Config, help string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var result Config
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
		10,
		"Number of winning spots to produce on program output.")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	if result.Multiplier < 1 || result.Multiplier > 10_000 {
		return nil, "", fmt.Errorf("Multiplier out of range 1 to 10000: %v", result.Multiplier)
	}
	result.args = flags.Args()
	if len(result.args) == 1 && len(result.args[0]) == 0 {
		// I love unit tests. Would have never figured it out without them.
		result.args = []string{}
	}

	result.SeedFromSystem = false
	// Loop here because Seed=0 has special meaning. We don't want it if system entropy
	// returns this value, but we also do not want infinite loop if the system is broken
	for i := 0; result.Seed == 0; i++ {
		if i > 3 {
			return nil, "", fmt.Errorf("System entropy failure, 0 returned multiple times")
		}

		const desired = 16
		b := make([]byte, desired)
		size, err := rand.Read(b)
		if err != nil {
			return nil, "", fmt.Errorf("Failed to read system entropy: %w", err)
		} else if size < desired {
			return nil, "", fmt.Errorf("Failed to read desired size from system entropy")
		}
		result.Seed = int64(binary.BigEndian.Uint64(b))
		result.SeedFromSystem = true
	}
	return &result, "", nil
}

func InputFile(config *Config) (file *os.File, err error) {
	size := len(config.args)
	if size == 0 {
		return os.Stdin, nil
	} else if size > 1 {
		return nil, fmt.Errorf("Only one filename expected, but found %v: %v", size, config.args)
	}

	fd, err := os.Open(config.args[0])
	if err != nil {
		return nil, fmt.Errorf("Unable to open input file: %w", err)
	}
	return fd, nil
}
