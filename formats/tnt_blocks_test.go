package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
)

func TestFourSpeciesWriteBlocks(t *testing.T) {
	var err error
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTACG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte("ATAGTCACG"))
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	sequence3 := sequence.NewSequence("Homo sapiens", []byte("TAGCATAGCTG"))
	sequence3.Species = "Homo sapiens"
	sequence3.Gene = "ATP6"

	sequence4 := sequence.NewSequence("Homo erectus", []byte("TAGCATAGCTA"))
	sequence4.Species = "Homo erectus"
	sequence4.Gene = "ATP6"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2, sequence3, sequence4)
	tnt.MetaData, err = tnt.GenerateMetaData()
	tnt.MetaData.Sort()

	if err != nil {
		t.Errorf("Error should be nil, was %s", err.Error())
	}

	buf := bytes.Buffer{}

	tnt.WriteBlocks(&buf)

	expected := `
blocks 0 11;
cnames
[1 ATP6;
[2 ATP8;
;`

	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}

}
