## Random selection of a raffle winner

The purpose is clear enough, as stated above.

The input is .csv file which lists entries and number of tickets per entry.

The basis of operation:
- read .csv file
  * two columns, first is name, second is numerical value
  * allow for duplicate names (i.e. someone bought tickets more than once)
  * allow for fractional number of "tickets", which will be rounded to a selected precision (e.g. currency can be used as "tickets")
  * report total number of entries and total number of tickets
- explode every entry into separate tickets
  * allow for a very large total number of tickets
- random shuffle all the tickets
  * random seed comes from the program input
- display top N entries
  * duplicate names are not removed
