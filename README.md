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
Day: Mon
Day: Mon
...
```

### Example: Piping and Aggregation

Count errors per day in an Apache log file:

```sh
./patt "[<day> <_>] [error] <_>" "Day: <day>" ./test_files/Apache_2k.log | sort | uniq -c | ./patt " <count> Day: <day>" "There were <count> errors on <day>"
There were    284 errors on Mon
There were    311 errors on Sun
```

## Benchmark

```sh
$ hyperfine "./pattern-cli '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log"     "awk '/\[error\]/ { if (match(\$0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print \"Day: \" m[1] }' ./test_files/Apache_2k.log"
hyperfine "./patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log"     "awk '/\[error\]/ { if (match(\$0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print \"Day: \" m[1] }' ./test_files/Apache_2k.log"
Benchmark 1: ./patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log
  Time (mean ± σ):       6.1 ms ±   0.8 ms    [User: 2.8 ms, System: 4.1 ms]
  Range (min … max):     5.0 ms …  10.6 ms    279 runs

  Warning: Command took less than 5 ms to complete. Results might be inaccurate.

Benchmark 2: awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print "Day: " m[1] }' ./test_files/Apache_2k.log
  Time (mean ± σ):      27.8 ms ±   5.8 ms    [User: 22.8 ms, System: 4.7 ms]
  Range (min … max):    13.4 ms …  53.8 ms    63 runs

  Warning: Statistical outliers were detected. Consider re-running this benchmark on a quiet PC without any interferences from other programs. It might help to use the '--warmup' or '--prepare' options.

Summary
  './patt '[<day> <_>] [error] <_>' 'Day: <day>' ./test_files/Apache_2k.log' ran
    4.55 ± 1.13 times faster than 'awk '/\[error\]/ { if (match($0, /^\[([A-Za-z]+) .*\] \[error\]/, m)) print "Day: " m[1] }' ./test_files/Apache_2k.log'
```

## Why not just use grep?

- For simple searching, `grep` is faster and more flexible.
- For extracting and reformatting structured data, `patt` is much more convenient and often faster than `awk` or custom scripts.

## Installation

```sh
go build -o patt ./cmd/cli/
```

## License

This project is based on code from Grafana Loki, which is licensed under the GNU Affero General Public License v3 (AGPLv3). See the LICENSE file for details.
