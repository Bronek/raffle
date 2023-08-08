# raffle

The purpose of this program is random selection of a raffle winner(s)

The input is .csv file, which lists entries and number of tickets per entry. Since raffle relies on random shuffle which needs to be properly initialised, seed for PRNG is very important, and can be provided as a program option. Same seed and input will result in the same draw results.

## The basis of operation:

Program performs following steps:

1. Read .csv file
  * two columns, first is name, second is numerical value
  * allow for duplicate names (i.e. someone bought tickets more than once)
  * allow for fractional number of "tickets", which will be rounded to a selected precision (e.g. currency can be used as "tickets")
  * report total number of tickets

2. Explode every entry into separate tickets
  * program can handle large total number of tickets

3. Random shuffle all the tickets
  * random seed can be provided in program options; if not provided the program will use own seed (derived from time)

4. Display top N entries
  * duplicate names are not removed, so someone who bought a large number of tickets can appear more than once on output.

## Program options

Program takes the following options:

`--input` File to read entries from. Must be in comma-separated-values format, where first column is name and second column is number of tickets (decimal point is allowed). If filename is not provided, program will read entries from standard input.

`--seed` Seed for PRNG. Same seed value will result in the same shuffle of the input data. If seed is not provided, program will use nanoseconds from the current time.

`--multiplier` Multiplier for the number of tickets. This is used to multiply input data, which is next rounded and used to generate tickets for each entry. For example if `--multiplier=100` and input contains number of tickets such as `1.20`, this will generate 120 tickets for such entry. Default 1, i.e. example `1.20` will generate 1 ticket.

`--N` How many lines to show on program output. Default 1, i.e. only one winner.

# TODO

- Improve error handling
- Add unit tests
- Improve default sorce of entropy
