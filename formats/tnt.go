package formats

import (
	"io"

	"github.com/alecthomas/template"
	"github.com/yarbelk/refasta/sequence"
)

// TNT formatter
type TNT struct {
	Title     string
	Sequences []sequence.Sequence
}

const tntNonInterleavedTemplateString = `'{{ .Title }}'
{{ .Length }} {{ .NTaxa }}
{{ range $i, $seq := .Sequences }}{{ $seq.SafeSpecies }} {{ $seq.Seq }}
{{ end }};`

var tntNonInterleavedTemplate = template.Must(template.New("TNT").Parse(tntNonInterleavedTemplateString))

// AddSequence to the internal sequence store.
func (t *TNT) AddSequence(seq sequence.Sequence) {
	t.Sequences = append(t.Sequences, seq)
}

// WriteSequences will collect up the sequences, verify their validity,
// and output a formated TNT file to the supplied writer
func (t *TNT) WriteSequences(writer io.Writer) error {
	if err := t.CheckSequenceLengths(); err != nil {
		return err
	}
	context := struct {
		Title         string
		Length, NTaxa int
		Sequences     []sequence.Sequence
	}{
		Title:     t.Title,
		Length:    t.Sequences[0].Length,
		NTaxa:     len(t.Sequences),
		Sequences: t.Sequences,
	}
	return tntNonInterleavedTemplate.Execute(writer, context)
}

// checkSequenceLengths will make sure that the sequences for the same
// gene sequence (or whatever sequence) are all the same length.
// Returns types of InvalidSequence with ErrNo
// MISSMATCHED_SEQUENCE_LENGTHS if they are no correct
func (t *TNT) CheckSequenceLengths() error {
	var expectedLen int
	for i, seq := range t.Sequences {
		if i == 0 {
			expectedLen = seq.Length
			continue
		}
		if seq.Length != expectedLen {
			return sequence.InvalidSequence{
				Message: "Sequences are not the Same length",
				Details: "None so far",
				Errno:   sequence.MISSMATCHED_SEQUENCE_LENGTHS,
			}
		}
	}
	return nil
}
