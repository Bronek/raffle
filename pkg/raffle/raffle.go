package raffle

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/bronek/raffle/pkg/options"
)

func Open() *os.File {
	if len(*options.Input) == 0 {
		return os.Stdin
	}

	fd, err := os.Open(*options.Input)
	if err != nil {
		panic("Unable to open input file")
	}
	return fd
}

func reader(file *os.File) *csv.Reader {
	ret := csv.NewReader(file)

	ret.FieldsPerRecord = 2
	ret.TrimLeadingSpace = true
	ret.Comment = '#'
	return ret
}

type Record struct {
	Name  string
	Count int
}

func Input(file *os.File) []Record {
	csvr := reader(file)
	result := make([]Record, 0, 0)

	for {
		row, err := csvr.Read()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			return result
		}

		count, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			panic(err)
		}
		count = math.Round(count * float64(*options.Multiplier))
		result = append(result, Record{row[0], int(count)})
	}
}

func Prepare(input []Record) []string {
	total := 0
	for i := len(input); i > 0; i-- {
		total += input[i-1].Count
	}

	fmt.Println("Total", float64(total)/float64(*options.Multiplier))
	result := make([]string, total, total)

	k := 0
	for i := len(input); i > 0; i-- {
		for j := 0; j < input[i-1].Count; j++ {
			result[k] = input[i-1].Name
			k++
		}
	}

	return result
}

func Shuffle(seed int64, exploded []string) []string {
	rand.Seed(seed)
	rand.Shuffle(len(exploded), func(i, j int) {
		exploded[i], exploded[j] = exploded[j], exploded[i]
	})
	return exploded
}
