package formats

import (
	"io"

	"github.com/alecthomas/template"
	"github.com/yarbelk/refasta/sequence"
)

type TNT struct {
	Title     string
	Sequences []sequence.Sequence
}

const tntNonInterleavedTemplateString = `'{{ .Title }}'
{{ .Length }} {{ .NTaxa }}
{{ range $i, $seq := .Sequences }}{{ $seq.SafeSpecies }} {{ $seq.Seq }}
{{ end }};`

var tntNonInterleavedTemplate = template.Must(template.New("TNT").Parse(tntNonInterleavedTemplateString))

func (t *TNT) AddSequence(seq sequence.Sequence) {
	t.Sequences = append(t.Sequences, seq)
}

func (t *TNT) WriteSequences(writer io.Writer) error {
	context := struct {
		Title         string
		Length, NTaxa int
		Sequences     []sequence.Sequence
	}{
		Title:     t.Title,
		Length:    len(t.Sequences[0].Seq),
		NTaxa:     len(t.Sequences),
		Sequences: t.Sequences,
	}
	return tntNonInterleavedTemplate.Execute(writer, context)
}
