package options

import (
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
)

var Input = flag.String(
	"input",
	"",
	"Input file, expected format CSV with two columns: name, number of tickets (may contain decimal point, see `multiplier` option). Defaults to standard input.")
var Seed = flag.Int64("seed", 0, "Seed value for random shuffle, must not be equal 0. If not set, will use system entropy")
var Multiplier = flag.Int("multiplier", 1, "Ticket size multiplier, must be greater than 0. Use 100 to handle ticket sizes with two decimal places (e.g. currency)")
var N = flag.Int("N", 1, "Number of winning spots to produce")

func init() {
	flag.Parse()

	if *Multiplier < 1 {
		panic("Invalid multiplier")
	}

	fmt.Println("Options")
	fmt.Println(" --input ", *Input)
	systemEntropy := false
	for *Seed == 0 {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		*Seed = int64(binary.BigEndian.Uint64(b))
		systemEntropy = true
	}
	if systemEntropy {
		fmt.Println(" --seed", *Seed, "(from system entropy)")
	} else {
		fmt.Println(" --seed", *Seed)
	}

	fmt.Println(" --multiplier", *Multiplier)
	fmt.Println(" --N", *N)
}
