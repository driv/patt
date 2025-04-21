package patt_test

import (
	"bytes"
	"os"
	"testing"

	patt "pattern"
)

func TestRunCLI_InputFromReader_OutputToWriter(t *testing.T) {
	input := "test line\nanother line\n"
	expectedOutput := "test line\n"
	pattern := "test<_>"

	// Create input and output buffers
	inputReader := bytes.NewReader([]byte(input))
	outputWriter := &bytes.Buffer{}

	// Run the CLI function
	err := patt.RunCLIWithIO(pattern, inputReader, outputWriter)
	if err != nil {
		t.Fatalf("RunCLIWithIO returned an error: %v", err)
	}

	// Verify the output
	if outputWriter.String() != expectedOutput {
		t.Errorf("unexpected output: got %q, want %q", outputWriter.String(), expectedOutput)
	}
}

func TestRunCLI_InvalidInput(t *testing.T) {
	invalidInput := ""
	expectedOutput := ""
	pattern := "test"

	// Create input and output buffers
	inputReader := bytes.NewReader([]byte(invalidInput))
	outputWriter := &bytes.Buffer{}

	// Run the CLI function
	err := patt.RunCLIWithIO(pattern, inputReader, outputWriter)
	if err == nil {
		t.Errorf("RunCLIWithIO should return an error: %v", err)
	}

	// Verify the output
	if outputWriter.String() != expectedOutput {
		t.Errorf("unexpected output: got %q, want %q", outputWriter.String(), expectedOutput)
	}
}

func TestRunCLI_InputFromFile_OutputToFile(t *testing.T) {
	inputFile := "test_input.txt"
	outputFile := "test_output.txt"
	pattern := "test<_>"

	// Create a temporary input file
	err := os.WriteFile(inputFile, []byte("test line\nanother line\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	// Run the CLI function
	err = patt.RunCLI(pattern, inputFile, outputFile)
	if err != nil {
		t.Fatalf("RunCLI returned an error: %v", err)
	}

	// Verify the output file
	output, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	expectedOutput := "test line\n"
	if string(output) != expectedOutput {
		t.Errorf("unexpected output: got %q, want %q", string(output), expectedOutput)
	}
}

func TestRunCLI_InvalidInputFile(t *testing.T) {
	invalidFile := "nonexistent.txt"
	outputFile := "test_output.txt"
	pattern := "test"

	defer os.Remove(outputFile)

	// Run the CLI function
	err := patt.RunCLI(pattern, invalidFile, outputFile)
	if err == nil {
		t.Errorf("expected error for invalid input file, but got none")
	}
}

func TestRunCLI_InvalidOutputFile(t *testing.T) {
	inputFile := "test_input.txt"
	invalidOutputFile := "/invalid/path/output.txt"
	pattern := "test"

	// Create a temporary input file
	err := os.WriteFile(inputFile, []byte("test line\nanother line\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}
	defer os.Remove(inputFile)

	// Run the CLI function
	err = patt.RunCLI(pattern, inputFile, invalidOutputFile)
	if err == nil {
		t.Errorf("expected error for invalid output file, but got none")
	}
}

func BenchmarkRunCLI(b *testing.B) {
	inputFile := "/tmp/log/syslog"
	outputFile := "/dev/null"
	pattern := "<_> pop-os <_>"

	for i := 0; i < b.N; i++ {
		err := patt.RunCLI(pattern, inputFile, outputFile)
		if err != nil {
			b.Fatalf("RunCLI returned an error: %v", err)
		}
	}
}
