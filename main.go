package main

import (
	"fmt"
	"io"
	"os"

	flag "github.com/ogier/pflag"
	"github.com/yarbelk/refasta/formats"
)

type FakeReadCloser struct {
	io.Reader
}

func (f FakeReadCloser) Close() error {
	return nil
}

type FakeWriteCloser struct {
	io.Writer
}

func (f FakeWriteCloser) Close() error {
	return nil
}

func getOutputFilePointer(filename string) (io.WriteCloser, error) {
	if filename == "--" {
		return FakeWriteCloser{os.Stdout}, nil
	}
	return os.Create(filename)
}

func getInputFilePointer(filename string) (io.ReadCloser, error) {
	if filename == "--" {
		return FakeReadCloser{os.Stdin}, nil
	}
	return os.Open(filename)
}

func main() {
	input := flag.StringP("input-file", "i", "--", "input file, it must be a valid input, or '--', or blank.  if blank. or '--', will read from stdin")
	output := flag.StringP("output-file", "o", "--", "output file, it must be a valid input, or '--', or blank.  if blank. or '--', will write to stdout")

	flag.Parse()

	fmt.Println("here we go")
	inputFP, err := getInputFilePointer(*input)
	defer inputFP.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Input file was not valid", err)
	}
	outputFP, err := getOutputFilePointer(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Output file was not valid", err)
	}
	defer outputFP.Close()

	fasta := formats.Fasta{}
	err = fasta.Parse(inputFP)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(fasta.Sequences)
	fasta.WriteSequences(outputFP)
}
