# Patt

A fast CLI tool for log pattern matching and replacement, inspired by [Grafana Loki's pattern queries](https://grafana.com/docs/loki/latest/query/log_queries/#pattern).

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
./patt '<pattern>' '<replacement>' <input-file>
```

- `<pattern>`: Loki-style pattern, e.g. `[<day> <_>] [error] <_>`
- `<replacement>`: Output template using named captures, e.g. `Day: <day>`
- `<input-file>`: Path to the log file (or use `-` for stdin)

### Example: Extracting Days from Error Logs

```sh
./patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log
Day: Jan
Day: Jan
...
```

### Example: Piping and Aggregation

Count errors per day in a huge log file:

```sh
./patt "[<day> <_>] [error] <_>" "Day: <day>" ./test_files/Apache_2k.log | sort | uniq -c | ./patt " <count> Day: <day>" "There were <count> errors on <day>"
There were    284 errors on Mon
There were    311 errors on Sun
```

## Benchmark

```sh
hyperfine "./patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log"     "awk '/[error]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print \"Day: \" m[1] }' ./test_files/Apache_2k.log"

Benchmark 1: ./patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log
  Time (mean ± σ):       4.7 ms ±   2.7 ms    [User: 2.5 ms, System: 3.0 ms]
  Range (min … max):     0.0 ms …   9.0 ms    324 runs

Benchmark 2: awk '/[error]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print "Day: " m[1] }' ./test_files/Apache_2k.log
  Time (mean ± σ):      25.2 ms ±   5.6 ms    [User: 21.7 ms, System: 3.4 ms]
  Range (min … max):     7.5 ms …  50.4 ms    95 runs

Summary
  './patt ...' ran 5.33 ± 3.25 times faster than the equivalent awk command.
```

## Why not just use grep?

- For simple searching, `grep` is faster and more flexible.
- For extracting and reformatting structured data, `patt` is much more convenient and often faster than `awk` or custom scripts.

## Installation

```sh
go build -o patt ./cmd/cli/main.go
```

## License

This project is based on code from Grafana Loki, which is licensed under the GNU Affero General Public License v3 (AGPLv3). See the LICENSE file for details.
