package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
)

func TestTwoSpiecesWithSameLengthData(t *testing.T) {
	sequence1 := sequence.Sequence{
		Species: "Homo sapiens",
		Gene:    "ATP8",
		Seq:     []byte("ATAGCTACG"),
	}
	sequence2 := sequence.Sequence{
		Species: "Homo erectus",
		Gene:    "ATP8",
		Seq:     []byte("ATAGTCACG"),
	}

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1)
	tnt.AddSequence(sequence2)

	buf := bytes.Buffer{}

	tnt.WriteSequences(&buf)

	expected := `'Title Here'
9 2
Homo_sapiens ATAGCTACG
Homo_erectus ATAGTCACG
;`
	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}
}
