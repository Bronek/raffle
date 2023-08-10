package options

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type Config struct {
	Seed            int64
	SeedFromEntropy bool
	Multiplier      int
	N               int
	args            []string
}

func (config Config) String() string {
	fromEntropySource := "set explicitly"
	if config.SeedFromEntropy {
		fromEntropySource = "entropy source"
	}
	return fmt.Sprintf("seed=%v (%s) multiplier=%v N=%v input=%+v",
		config.Seed, fromEntropySource, config.Multiplier, config.N, config.args)
}

func Parse(progname string, args []string, entropy io.Reader) (config *Config, help string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var result Config
	flags.Int64Var(
		&result.Seed,
		"seed",
		0,
		"Random seed for shuffle operation. If not provided or set to 0,\n"+
			"program will use system entropy source")
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

	result.SeedFromEntropy = false
	// Loop here because Seed=0 has special meaning. We don't want it even if entropy source
	// returns this value, but we also do not want infinite loop if the entropy source is broken.
	// Retries allow us to use real-world entropy sources, which can be intermittent or slow.
	desired := 8
	retries := 3
	random := make([]byte, 0, desired)
	for i := 0; result.Seed == 0; i++ {
		if i > retries {
			return nil, "", fmt.Errorf("Entropy source failure, 0 returned repeatedly")
		}
		// Real-world entropy source may need time to collect entropy
		time.Sleep(time.Duration(i) * time.Millisecond)

		b := make([]byte, desired)
		size, err := entropy.Read(b)
		if err != nil {
			if i == retries {
				return nil, "", fmt.Errorf("Failed to read entropy source: %w", err)
			} else {
				continue
			}
		}

		// Filter out all all-zeros-reads from entropy source.
		allzeros := true
		for j := 0; j < size; j++ {
			if b[j] != 0 {
				allzeros = false
			}
		}
		if allzeros {
			continue
		}

		random = append(random, b[:size]...)
		desired -= size
		if desired == 0 {
			result.Seed = int64(binary.BigEndian.Uint64(random))
			result.SeedFromEntropy = true
			// Since we filtered out allzeros, we know at this point that Seed != 0
		} else if i == retries {
			return nil, "", fmt.Errorf("Failed to read desired size from entropy source")
		}
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
