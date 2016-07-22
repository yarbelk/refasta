package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
	"gopkg.in/urfave/cli.v1"
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
	if err != nil {
		return false, err
	}
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

	for _, file := range files {
		ext := path.Ext(file)
		geneName := filepath.Base(file[:len(file)-len(ext)])
		err := func() error {
			fasta := formats.Fasta{SpeciesFromID: true}
			fd, err := os.Open(file)
			if err != nil {
				// probably an Access Control issue, or race condition
				return err
			}
			defer fd.Close()
			err = fasta.Parse(fd, geneName)
			if err != nil {
				// Some parsing error...
				return err
			}
			sequences = append(sequences, fasta.Sequences...)
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return sequences, nil
}

func handleFastaOutput(sequences []sequence.Sequence, output string) error {
	fasta := formats.Fasta{}
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

	app := cli.NewApp()
	app.Name = "refasta"
	app.Usage = `Convert various genitics data formats into other formats.
	Currently only fasta and tnt are supported, and in an opinionated way.`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input-format, f",
			Value: formats.FASTA_FORMAT,
			Usage: "INPUT_FORMAT must be one of the supported input types. Currently supported is 'fasta'",
		},
		cli.StringFlag{
			Name:  "output-format, F",
			Value: formats.FASTA_FORMAT,
			Usage: "OUTPUT_FORMAT must be one of the supported output types. Currently supported are 'fasta' and 'tnt'",
		},
		cli.StringFlag{
			Name:  "input-file, i",
			Value: "--",
			Usage: "INPUT_FILE, it must be a valid input, or '--', or blank.  if blank. or '--', will read from stdin",
		},
		cli.StringFlag{
			Name:  "output-file, o",
			Value: "--",
			Usage: "OUTPUT_FILE, it must be a valid input, or '--', or blank.  if blank. or '--', will write to stdout",
		},
		cli.StringFlag{
			Name:  "tnt-title, t",
			Value: "",
			Usage: "title for TNT output",
		},
		cli.StringFlag{
			Name:  "outgroup",
			Value: "",
			Usage: "Optional OUTGROUP for TNT output.  If specified, this species will be used as the outgroup for TNT. " +
				"Otherwise the first (alphabetically) will be used.  This must be left blank, or be a valid species name " +
				"from the input",
		},
	}
	fatalWithUsage := func(c *cli.Context, err error) {
		if err != nil {
			fmt.Fprintln(c.App.Writer, err.Error())
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
	}

	app.Action = func(c *cli.Context) {
		var sequences []sequence.Sequence
		var err error

		input := c.String("input-file")
		output := c.String("output-file")
		inputFormat := c.String("input-format")
		outputFormat := c.String("output-format")

		switch inputFormat {
		case formats.FASTA_FORMAT:
			sequences, err = handleFastaInput(input)
		default:
			err = fmt.Errorf("Unknown intput format '%s'", inputFormat)
		}
		fatalWithUsage(c, err)

		switch outputFormat {
		case formats.FASTA_FORMAT:
			err = handleFastaOutput(sequences, output)
		case formats.TNT_FORMAT:
			fmt.Fprintf(os.Stderr, "Output format is TNT; serializing\n")
			context := TNTContext{
				Title:    c.String("title"),
				Outgroup: c.String("outgroup"),
			}
			err = handleTNTOutput(context, sequences, output)
		default:
			err = fmt.Errorf("Unknown output format '%s'", inputFormat)
		}
		fatalWithUsage(c, err)
	}

	app.Run(os.Args)
}
