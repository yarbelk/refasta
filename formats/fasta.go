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
type FastaWriter struct {
	FileName  string
	File      io.Writer
	sequences []sequence.Sequence
}

// AddSequence to the internal list of sequences of the writer
func (fw *FastaWriter) AddSequence(seq sequence.Sequence) {
	fw.sequences = append(fw.sequences, seq)
}

// WriteSequences writes the stored sequences to the stored file pointer
func (fw *FastaWriter) WriteSequences() error {
	return fastaTemplate.Execute(fw.File, fw.sequences)
}
