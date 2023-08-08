package options

import (
	"flag"
	"fmt"
	"time"
)

var Input = flag.String(
	"input",
	"",
	"Input file, expected format CSV with two columns: name, number of tickets (may contain decimal point, see `multiplier` option). Defaults to standard input.")
var Seed = flag.Int64("seed", 0, "Seed value for random shuffle, must be greater than 0. If not set, will use current time")
var Multiplier = flag.Int("multiplier", 1, "Ticket size multiplier, must be greater than 0. Use 100 to handle ticket sizes with two decimal places (e.g. currency)")
var Unique = flag.Bool("unique", false, "Unique winners. This prevents single winner from taking more than one winning spot, even if they have most of the tickets")
var N = flag.Int("N", 1, "Number of winning spots to produce")

func init() {
	flag.Parse()

	if *Multiplier < 1 {
		panic("Invalid multiplier")
	}

	fmt.Println("Options")
	fmt.Println(" --input ", *Input)
	if *Seed < 1 {
		*Seed = (time.Now().UnixNano() % 1000000000)
		fmt.Println(" --seed", *Seed, "(from time)")
	} else {
		fmt.Println(" --seed", *Seed)
	}
	fmt.Println(" --multiplier", *Multiplier)
	fmt.Println(" --unique", *Unique)
	fmt.Println(" --N", *N)
}
