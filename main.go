package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	flag "github.com/ogier/pflag"
	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
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

type TNTContext struct {
	Title    string
	Outgroup string
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

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func isFormat(file, format string) bool {
	switch path.Ext(file) {
	case ".fas", ".fasta":
		return format == formats.FASTA_FORMAT
	default:
		return false
	}
}

func dirInput(dir, format string, recurse bool) ([]string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	filteredFiles := make([]string, 0, 10)
	for _, file := range files {
		if file.IsDir() && recurse {
			f, err := dirInput(filepath.Join(absDir, file.Name()), format, true)
			if err != nil {
				return nil, err
			}
			filteredFiles = append(filteredFiles, f...)
		} else if isFormat(file.Name(), format) {
			filteredFiles = append(filteredFiles, filepath.Join(absDir, file.Name()))
		}
	}
	return filteredFiles, nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func handleFastaInput(input string) ([]sequence.Sequence, error) {
	var files []string
	var sequences []sequence.Sequence

	if isDir, err := isDirectory(input); isDir && err == nil {
		files, err = dirInput(input, formats.FASTA_FORMAT, true)
		if err != nil {
			// Some error in walking the directory tree
			return nil, err
		}
	} else {
		files = []string{input}
	}
	fmt.Println("files", files)

	for _, file := range files {
		ext := path.Ext(file)
		geneName := file[:len(file)-len(ext)]
		fasta := formats.Fasta{SpeciesFromID: true}
		fd, err := os.Open(file)
		if err != nil {
			// probably an Access Control issue, or race condition
			return nil, err
		}
		defer fd.Close()
		err = fasta.Parse(fd, geneName)
		if err != nil {
			// Some parsing error...
			return nil, err
		}
		sequences = append(sequences, fasta.Sequences...)
	}
	return sequences, nil
}

func handleFastaOutput(sequences []sequence.Sequence, output string) error {
	fasta := formats.Fasta{}
	fmt.Println(sequences)
	fasta.AddSequence(sequences...)
	fd, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue opening output file;\n%s", err.Error())
	}
	defer fd.Close()
	return fasta.WriteSequences(fd)
}

func handleTNTOutput(context TNTContext, sequences []sequence.Sequence, output string) error {
	tnt := formats.TNT{Title: context.Title}
	fmt.Println(sequences)
	tnt.AddSequence(sequences...)
	tnt.SetOutgroup(context.Outgroup)
	fd, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue opening output file;\n%s", err.Error())
	}
	defer fd.Close()
	return tnt.WriteSequences(fd)
}

func main() {
	inputFormat := flag.StringP("input-format", "f", formats.FASTA_FORMAT, "intput format, must be 'fasta'")
	outputFormat := flag.StringP("outpu-format", "F", formats.TNT_FORMAT, "output format, must be [fasta|tnt]")
	input := flag.StringP("input-file", "i", "--", "input file, it must be a valid input, or '--', or blank.  if blank. or '--', will read from stdin")
	output := flag.StringP("output-file", "o", "--", "output file, it must be a valid input, or '--', or blank.  if blank. or '--', will write to stdout")
	tntTitle := flag.StringP("tnt-title", "t", "", "title for TNT output")
	outgroup := flag.String("outgroup", "", "outgroup for TNT")
	flag.Parse()

	var sequences []sequence.Sequence
	var err error

	switch *inputFormat {
	case formats.FASTA_FORMAT:
		fmt.Fprintf(os.Stderr, "intput format is fasta; parsing\n")
		sequences, err = handleFastaInput(*input)
	default:
		fmt.Fprintf(os.Stderr, "Unknown intput format '%s'", inputFormat)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	switch *outputFormat {
	case formats.FASTA_FORMAT:
		fmt.Fprintf(os.Stderr, "Output format is fasta; serializing\n")
		if err := handleFastaOutput(sequences, *output); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	case formats.TNT_FORMAT:
		fmt.Fprintf(os.Stderr, "Output format is TNT; serializing\n")
		context := TNTContext{
			Title:    *tntTitle,
			Outgroup: *outgroup,
		}
		if err := handleTNTOutput(context, sequences, *output); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown output format '%s'", inputFormat)
		os.Exit(1)
	}

}
