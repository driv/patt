# Pattern Matching CLI Tool

## Overview

This repository contains a CLI tool for pattern matching in text files. The tool reads input from a file or standard input, matches lines against a specified pattern, and writes the matching lines to standard output. If no matches are found, an appropriate message is displayed.

## Features

- Supports pattern matching using a custom syntax.
- Reads input from a file or standard input.
- Outputs matching lines to standard output.
- Provides clear error messages for invalid patterns or input issues.

## File Structure

- `cmd/cli/main.go`: Entry point for the CLI tool.
- `pattern.go`: Contains the core pattern matching logic.
- `filter_lines.go`: Handles filtering and writing lines based on patterns.
- `pattern_test.go`: Unit tests for the pattern matching logic.
- `cli_test.go`: Tests for the CLI tool.
- `go.mod` and `go.sum`: Go module files.

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd patt
   ```

2. Build the CLI tool:

   ```bash
   go build -o pattern-cli ./cmd/cli
   ```

## Usage

Run the CLI tool with the following syntax:

```bash
./pattern-cli <pattern> [input-file]
```

- `<pattern>`: The pattern to match lines against.
- `[input-file]`: (Optional) Path to the input file. If omitted, the tool reads from standard input.

### Examples

1. Match lines containing the word "error" in a file:

   ```bash
   ./pattern-cli "error" log.txt
   ```

2. Match lines from standard input:

   ```bash
   echo -e "line1\nerror line\nline3" | ./pattern-cli "error"
   ```

## Testing

Run the tests using:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a clear description of your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
