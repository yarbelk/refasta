package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
)

// TestTwoGenesTwoSpecies should sort the blocks by alphabetical order
func TestTwoGenesTwoSpeciesFullOutput(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTAG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte("ATAGCTAC"))
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	sequence3 := sequence.NewSequence("Homo sapiens", []byte("TAGCATAGCTG"))
	sequence3.Species = "Homo sapiens"
	sequence3.Gene = "ATP6"

	sequence4 := sequence.NewSequence("Homo erectus", []byte("TAGCATAGCTA"))
	sequence4.Species = "Homo erectus"
	sequence4.Gene = "ATP6"

	expected := `xread
'Title Here'
19 2
Homo_erectus TAGCATAGCTAATAGCTAC
Homo_sapiens TAGCATAGCTGATAGCTAG
;
blocks 0 11;
cnames
[1 ATP6;
[2 ATP8;
;`

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2, sequence3, sequence4)

	buf := bytes.Buffer{}

	if err := tnt.WriteSequences(&buf); err != nil {
		t.Error("Expected no error, got one", err)
	}

	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}

}
