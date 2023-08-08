package main

import (
	"fmt"
	"math"

	"github.com/bronek/raffle/pkg/options"
	"github.com/bronek/raffle/pkg/raffle"
)

func main() {
	input := raffle.Input(raffle.Open())

	exploded := raffle.Shuffle(*options.Seed, raffle.Prepare(input))
	fmt.Println("Results")
	for i := 0; i < int(math.Min(float64(len(exploded)), float64(*options.N))); i++ {
		fmt.Println(exploded[i])
	}
}
