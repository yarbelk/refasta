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

var sequences []sequence.Sequence

type CommandError struct {
	error
	c *cli.Context
}

func (c CommandError) UsageAndFail() {
	if c.error != nil {
		fmt.Fprintln(c.c.App.Writer, c.error.Error())
		cli.ShowAppHelp(c.c)
		os.Exit(1)
	}

}

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

func parseInput(c *cli.Context) error {
	var err error
	var inputFormat string = c.String("intput-format")
	switch inputFormat {
	case formats.FASTA_FORMAT:
		sequences, err = handleFastaInput(c.String("input"))
	default:
		err = CommandError{fmt.Errorf("Unknown intput format '%s'", inputFormat), c}
	}
	return err
}

func main() {

	app := cli.NewApp()
	app.Name = "refasta"
	app.Usage = "Convert various genitics data formats into other formats. " +
		"Currently only fasta and tnt are supported, and in an opinionated way."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "",
			Usage: "`INPUT`, it must be either a file or directory.  If blank, stdin will be used",
		},
		cli.StringFlag{
			Name:  "input-format, f",
			Value: formats.FASTA_FORMAT,
			Usage: "`INPUT_FORMAT` must be one of the supported input types. Currently only 'fasta' is supported",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "fasta",
			Usage:       "Convert to `fasta` format",
			UsageText:   "This will convert the input to a fasta formatted file.",
			Description: "This requires an input file or directory, and an input format.  If you do not specify an OUTPUT_FILE, then the output will be written to stdout",
			ArgsUsage:   "[OUTPUT_FILE]",
			Before:      parseInput,
			Action: func(c *cli.Context) error {
				return handleFastaOutput(sequences, c.Args().First())
			},
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "tnt",
			Usage:       "Convert to `TNT` format",
			UsageText:   "This will convert the input to a TNT formatted file.",
			Description: "This requires an input file or directory, and an input format.  You can specify the outgroup and title of the file.  If you do not specify an OUTPUT_FILE, then the output will be written to stdout",
			ArgsUsage:   "[OUTPUT_FILE]",
			Before:      parseInput,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "outgroup",
					Value: "",
					Usage: "Optional `OUTGROUP` for TNT output.  If specified, this species will be used as the outgroup for TNT. " +
						"Otherwise the first (alphabetically) will be used.  This must be left blank, or be a valid species name " +
						"from the input",
				},
				cli.StringFlag{
					Name:  "tnt-title, t",
					Value: "",
					Usage: "`TITLE` for TNT output",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Fprintf(os.Stderr, "Output format is TNT; serializing\n")
				context := TNTContext{
					Title:    c.String("title"),
					Outgroup: c.String("outgroup"),
				}
				return handleTNTOutput(context, sequences, c.Args().First())
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		switch e := err.(type) {
		case CommandError:
			e.UsageAndFail()
		default:
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
