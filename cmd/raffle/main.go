package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/bronek/raffle/pkg/options"
	"github.com/bronek/raffle/pkg/raffle"
)

func main() {
	conf, help, err := options.Parse(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(help)
		os.Exit(2)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inputFile, err := options.InputFile(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	defer inputFile.Close()

	fmt.Printf("Configuration: {%v}\n", conf)

	input := raffle.Input(inputFile, conf.Multiplier)

	exploded := raffle.Shuffle(conf.Seed, raffle.Prepare(input, conf.Multiplier))
	fmt.Println("Results")
	for i := 0; i < int(math.Min(float64(len(exploded)), float64(conf.N))); i++ {
		fmt.Println(exploded[i])
	}
}
