# raffle

The purpose of this program is random selection of a raffle winner(s)

The input is provided in CSV format and contains a list of entries. Each entry consists of the
name of entrant, followed by comma and the number of tickets. Same entrant can be listed multiple
times. Since raffle relies on random shuffle, which needs to be properly initialised, seed can be
provided as a program option. Same seed and input will always result in the same draw results.

## The basis of operation:

Program performs following steps:

1. Read input file (if filename is not provided, program will read standard input)
  * two columns: first is name, second is number of tickets, separated by comma
  * allow for repeat names (e.g. entrant purchased tickets multiple times)
  * allow for fractional number of tickets. This will be first multiplied by a given value
    then rounded (e.g. using multiplier=100 allows using currency as "tickets")

2. Expand every entry into separate tickets
  * e.g. using multiplier=100 and number of tickets 2.5 will create 250 tickets for such entry
  * program can handle large total number of tickets

3. Random shuffle all the tickets
  * seed for shuffle can be provided in program options
  * program will use system entropy if seed is not provided
  * user can request system entropy explicitly by using seed=0
  * shuffle algorithm is Fisher-Yates (i.e. standard Go `math/rand.Shuffle`)

4. Display top N entries
  * duplicate names are not removed, so one entrant can appear more than once on program output

## Program options

Program takes the following options:

`-seed` Seed for PRNG. Same seed value will result in the same shuffle of the input data. If
seed is not provided or is explicitly set to 0, program will use system entropy.

`-multiplier` Multiplier for the number of tickets, must not exceed 10000. Default is 1.

`-N` How many winners to show on program output. Default is 10.

Options can be followed by name of the .csv file to read entries from. If filename is not provided,
program will read entries from the standard input.

## Building raffle binary

Run `go build ./cmd/raffle` in the project directory. Alternatively you can use the default target
of the convenience `Makefile`, which will also run tests and save the binary in `dist` subdirectory.

# Known issues

When building on Windows with go version 1.20, Windows Defender may break the build and report false
positive like one below. This is known problem with Go on Windows, explained here
https://go.dev/doc/faq#virus

    Operation did not complete successfully because the file contains a virus or potentially unwanted software.

Suggested workaround is to try a different go version on Windows; 1.21 and 1.19 are both expected
to work.

# TODO

- Improve error handling
- Add unit tests
