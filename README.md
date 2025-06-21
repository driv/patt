# Patt

A fast CLI tool for log pattern matching and replacement, based on [Grafana Loki's pattern queries](https://grafana.com/docs/loki/latest/query/log_queries/#pattern).

## Features

- **Pattern-based log matching** using a syntax similar to Loki's pattern queries.
- **Powerful replacement**: Extract and reformat log fields with named captures.
- **Streaming and piping**: Works well in Unix pipelines for further processing.
- **Performance**: Extremely fast for replacement tasks (see benchmarks below).

## When to Use

- **Best for**: Extracting and reformatting structured log data using named patterns.
- **Not ideal for**: Pure search/filtering (use `grep` for that).

## Usage

```sh
patt "<pattern>" ["<replacement>"] [-f <input-file>] [-k]
```

- `<pattern>`: Loki-style pattern, e.g. `[<day> <_>] [error] <_>`
- `<replacement>`: (Optional) Output template using named captures, e.g. `Day: <day>`
- `-f <input-file>`: (Optional, defaults to stdin) Path to the log file
- `-k, --keep`: (Optional) Print non-matching lines as well (like `sed`)

### Examples

#### Search Only

```sh
patt '[<day> <_>] [error] <_>' -f ./testdata/Apache_2k.log
```

- Prints lines matching the pattern from the file.

#### Replace (Extract and Reformat)

```sh
patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log
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
sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log
```

#### Piping and Aggregation

Count errors per day in an Apache log file:

```sh
patt "[<day> <_>] [error] <_>" "Day: <day>" -f ./testdata/Apache_2k.log | \
 sort | uniq -c | \
 patt " <count> Day: <day>" "There were <count> errors on <day>"

There were    284 errors on Mon
There were    311 errors on Sun
```

#### Highlight error lines, but keep all lines

```sh
patt '[<day> <_>] [error] <message>' '[ERROR] <day>: <message>' -f ./testdata/Apache_2k.log -k
```

- All lines are printed. Lines matching the pattern are reformatted to highlight errors; non-matching lines are printed unchanged.

## Benchmark

```sh
# Compared to awk

$ hyperfine "patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log"     "awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print \"Day: \" m[1] }' ./testdata/Apache_2k.log"

Benchmark 1: patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log
  Time (mean ± σ):       6.1 ms ±   0.8 ms    [User: 2.8 ms, System: 4.1 ms]
  Range (min … max):     5.0 ms …  10.6 ms    279 runs

  Warning: Command took less than 5 ms to complete. Results might be inaccurate.

Benchmark 2: awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print "Day: " m[1] }' ./testdata/Apache_2k.log
  Time (mean ± σ):      27.8 ms ±   5.8 ms    [User: 22.8 ms, System: 4.7 ms]
  Range (min … max):    13.4 ms …  53.8 ms    63 runs

Summary
  'patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log' ran
    4.55 ± 1.13 times faster than 'awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print "Day: " m[1] }' ./testdata/Apache_2k.log'
```

```sh
# Compared to sed

$ hyperfine "patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log"     "sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log"
Benchmark 1: patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log
  Time (mean ± σ):       5.5 ms ±   3.5 ms    [User: 2.1 ms, System: 3.7 ms]
  Range (min … max):     1.0 ms …  14.6 ms    402 runs

  Warning: Command took less than 5 ms to complete. Results might be inaccurate.

Benchmark 2: sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log
  Time (mean ± σ):      19.6 ms ±  12.2 ms    [User: 16.6 ms, System: 2.7 ms]
  Range (min … max):     7.1 ms …  44.9 ms    84 runs

Summary
  'patt '[<day> <_>] [error] <_>' 'Day: <day>' -f ./testdata/Apache_2k.log' ran
    3.54 ± 3.15 times faster than 'sed -n 's/^\[\([A-Za-z]\+\) .*\] \[error\].*/Day: \1/p' ./testdata/Apache_2k.log'
```

## Installation

```sh
go build -o patt ./cmd/patt/
```

## License

This project is based on code from Grafana Loki, which is licensed under the GNU Affero General Public License v3 (AGPLv3). See the LICENSE file for details.
