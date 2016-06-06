package formats

import (
	"fmt"
	"io"
	"text/template"

	"github.com/yarbelk/refasta/sequence"
)

const fastaTemplateString = `{{ range $i, $seq := . -}}
>{{- $seq.SafeName }}
{{ printf "%s" .Seq }}
{{ end }}
`

const FASTA_FORMAT = "fasta"

var fastaTemplate = template.Must(template.New("fasta").Parse(fastaTemplateString))

// FastaWriter writes a seriese of sequences to a fasta file
type Fasta struct {
	Sequences     []sequence.Sequence
	SpeciesFromID bool
}

// AddSequence (or many) to the internal list of sequences of the writer
func (f *Fasta) AddSequence(seqs ...sequence.Sequence) {
	f.Sequences = append(f.Sequences, seqs...)
}

// WriteSequences writes the stored sequences to the stored file pointer
func (f *Fasta) WriteSequences(writer io.Writer) error {
	return fastaTemplate.Execute(writer, f.Sequences)
}

// Parse will read a file, and append all new Sequences to the store
// of sequences
func (f *Fasta) Parse(input io.Reader, geneName ...string) error {
	var gene string
	if len(geneName) == 1 {
		gene = geneName[0]
	}
	fastaScanner := NewFastaScanner(input)
	var newSequence sequence.Sequence = sequence.NewSequence("", []byte{})
	var lastToken Token = UNSTARTED
	for {
		token, lit, alpha, length := fastaScanner.Scan()

		switch token {
		case SEQUENCE_ID:
			if lastToken == SEQUENCE_ID {
				return sequence.FormatError{
					Message: "Badly formated FASTA file",
					Details: "Two sequence id ('>....', without any data in between",
					Errno:   sequence.BAD_FORMAT,
				}
			}
			newSequence = sequence.Sequence{Name: string(lit), Gene: gene}
			(&newSequence).SetAlphabet(alpha)
			if f.SpeciesFromID {
				newSequence.Species = string(lit)
			}
			lastToken = SEQUENCE_ID
			continue
		case SEQUENCE_DATA:
			if lastToken != SEQUENCE_ID {
				return sequence.FormatError{
					Message: "Badly formated FASTA file",
					Details: "Sequence data did not have a Sequence ID",
					Errno:   sequence.BAD_FORMAT,
				}
			}
			newSequence.Seq = lit
			newSequence.Length = length
			lastToken = SEQUENCE_DATA
			f.Sequences = append(f.Sequences, newSequence)
		case EOF:
			return nil
		case INVALID:
			return fmt.Errorf("Invalid characters in the stream")
		}
	}
}
