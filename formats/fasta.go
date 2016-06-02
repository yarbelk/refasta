package formats

import (
	"io"
	"text/template"

	"github.com/yarbelk/refasta/sequence"
)

const fastaTemplateString = `{{ range $i, $seq := . -}}
>{{- $seq.SafeName }}
{{ printf "%s" .Seq }}
{{ end }}
`

var fastaTemplate = template.Must(template.New("fasta").Parse(fastaTemplateString))

// FastaWriter writes a seriese of sequences to a fasta file
type Fasta struct {
	Sequences []sequence.Sequence
}

// AddSequence to the internal list of sequences of the writer
func (f *Fasta) AddSequence(seq sequence.Sequence) {
	f.Sequences = append(f.Sequences, seq)
}

// WriteSequences writes the stored sequences to the stored file pointer
func (f *Fasta) WriteSequences(writer io.Writer) error {
	return fastaTemplate.Execute(writer, f.Sequences)
}

// Parse will read a file, and append all new Sequences to the store
// of sequences
func (f *Fasta) Parse(input io.Reader) error {
	return nil
}
