# Patt

A fast CLI tool for log pattern matching and replacement, based on [Grafana Loki's pattern queries](https://grafana.com/docs/loki/latest/query/log_queries/#pattern).

## Features

- **Pattern-based log matching** using a syntax similar to Loki's pattern queries.
- **Powerful replacement**: Extract and reformat log fields with named captures.
- **Multiple patterns and files**: Search for multiple patterns across multiple files.
- **Streaming and piping**: Works well in Unix pipelines for further processing.
- **Performance**: Extremely fast for replacement tasks (see benchmarks below).

## When to Use

- **Best for**: Extracting and reformatting structured log data using named patterns.
- **Not ideal for**: Pure search/filtering (use `grep` for that).

## Usage

```sh
patt [flags] <search_pattern> [<search_pattern>...] [<replacement>] [-- <input_file> [<input_file>...]]

```

- `<search_pattern>`: One or more Loki-style patterns, e.g. `[<day> <_>] [error] <_>`.
- `<replacement>`: (Optional) Output template using named captures, e.g. `Day: <day>`.
- `<input_file>`: (Optional, defaults to stdin) One or more paths to log files. Use `--` to separate files from patterns.
- `-k, --keep`: (Optional) Print non-matching lines as well (like `sed`).

### Examples

#### Search Only

```sh
patt '[<day> <_>] [error] <_>' -- ./testdata/Apache_2k.log
```

- Prints lines matching the pattern from the file.

#### Multiple Search Patterns

```sh
patt "[Sun Dec 04 <_>] <something>" "[Mon Dec 05 <_>] <something>" "Found on Sun or Mon: <something>" -- ./testdata/Apache_2k.log
```

- Prints lines matching either of the search patterns, formatted with the replacement.

#### Replace (Extract and Reformat)

```sh
patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_2k.log
Day: Mon
Day: Mon
<...>
```

#### From stdin

```sh
cat ./testdata/Apache_2k.log | patt '[<day> <_>] [error] <_>' 'Day: <day>'
```

#### Equivalent with sed

```sh
sed -n 's/^[A-Za-z]\+ .*] [error].*/Day: \1/p' ./testdata/Apache_2k.log
```

#### Piping and Aggregation

Count errors per day in an Apache log file:

```sh
patt "[<day> <_>] [error] <_>" "Day: <day>" -- ./testdata/Apache_2k.log | \
 sort | uniq -c | \
 patt " <count> Day: <day>" "There were <count> errors on <day>"

There were    284 errors on Mon
There were    311 errors on Sun
```

#### Highlight error lines, but keep all lines

```sh
patt '[<day> <_>] [error] <message>' '[ERROR] <day>: <message>' -k -- ./testdata/Apache_2k.log
```

- All lines are printed.

## Benchmark

```sh
# Compared to awk

$ hyperfine "./patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_500MB.log" "awk '/\[error\]/ { if (match(\$0, /^\[([A-Za-z]+)/, m)) print \"Day:\", m[1] }' ./testdata/Apache_500MB.log"
Benchmark 1: ./patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_500MB.log
  Time (mean ± σ):      1.085 s ±  0.121 s    [User: 1.006 s, System: 0.090 s]
  Range (min … max):    0.964 s …  1.359 s    10 runs
 
Benchmark 2: awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+)/, m)) print "Day:", m[1] }' ./testdata/Apache_500MB.log
  Time (mean ± σ):      6.727 s ±  0.600 s    [User: 6.571 s, System: 0.138 s]
  Range (min … max):    5.876 s …  7.645 s    10 runs
 
Summary
  './patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_500MB.log' ran
    6.20 ± 0.88 times faster than 'awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+)/, m)) print "Day:", m[1] }' ./testdata/Apache_500MB.log'
```

```sh
# Compared to sed

$ hyperfine "patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_2k.log" "sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log"
Benchmark 1: patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_2k.log
  Time (mean ± σ):       5.5 ms ±   3.5 ms    [User: 2.1 ms, System: 3.7 ms]
  Range (min … max):     1.0 ms …  14.6 ms    402 runs

  Warning: Command took less than 5 ms to complete. Results might be inaccurate.

Benchmark 2: sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log
  Time (mean ± σ):      19.6 ms ±  12.2 ms    [User: 16.6 ms, System: 2.7 ms]
  Range (min … max):     7.1 ms …  44.9 ms    84 runs

Summary
  'patt '[<day> <_>] [error] <_>' 'Day: <day>' -- ./testdata/Apache_2k.log' ran
    3.54 ± 3.15 times faster than 'sed -n 's/^[A-Za-z]\+ .*] [error].*/Day: \1/p' ./testdata/Apache_2k.log'
```

## Installation

```sh
go build -o patt ./cmd/patt/
```

## License

This project is based on code from Grafana Loki, which is licensed under the GNU Affero General Public License v3 (AGPLv3). See the LICENSE file for details.
