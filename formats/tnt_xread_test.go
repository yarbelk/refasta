package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
)

func TestTwoSpiecesWithSameLengthData(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTACG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte("ATAGTCACG"))
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2)

	buf := bytes.Buffer{}

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Error("Expected no error, got one", err)
	}

	tnt.WriteXRead(&buf)

	expected := `xread
'Title Here'
9 2
Homo_erectus ATAGTCACG
Homo_sapiens ATAGCTACG
;`
	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}
}

func TestCleanDataWorksWithTwoSpeciesWithOneMissingData(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTACG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte{})
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2)

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Errorf("Expected No Error with GenerateMetaData, got %v", err)
	}

	tnt.CleanData()
	expected := "---------"
	got := string(tnt.Sequences["ATP8"]["Homo erectus"].Seq)
	if got != expected {
		t.Errorf("Expected CleanData to fill in missing data for 'Homo erectus'. expected '%s', got '%s'",
			expected, got)
	}
}

func TestTwoSpeciesTwoGenesWithOneMissingDataFillsIn(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTACG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte{})
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	sequence3 := sequence.NewSequence("Homo sapiens", []byte("ATAGCACACTACG"))
	sequence3.Species = "Homo sapiens"
	sequence3.Gene = "ATP6"

	sequence4 := sequence.NewSequence("Homo erectus", []byte("ATAGCACACTACG"))
	sequence4.Species = "Homo erectus"
	sequence4.Gene = "ATP6"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2, sequence3, sequence4)

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Errorf("Expected No Error with GenerateMetaData, got %v", err)
	}

	buf := bytes.Buffer{}
	tnt.GenerateMetaData()
	tnt.CleanData()
	tnt.WriteXRead(&buf)

	expected := `xread
'Title Here'
22 2
Homo_erectus ATAGCACACTACG---------
Homo_sapiens ATAGCACACTACGATAGCTACG
;`

	got := buf.String()
	if got != expected {
		t.Errorf("Expected CleanData to fill in missing data for 'Homo erectus'. expected '%s', got '%s'",
			expected, got)
	}
}
func TestTwoSpeciesWithDifferingLengthDataHasError(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCTACG"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte("ATAGCTAC"))
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2)

	_, err := tnt.GenerateMetaData()

	if err == nil {
		t.Errorf("Expected tnt.WriteSequence to fail with a differing length issue")
	}

	switch is := err.(type) {
	case sequence.InvalidSequence:
		if is.Errno != sequence.MISSMATCHED_SEQUENCE_LENGTHS {
			t.Errorf("Expected is to be a 'MISSMATCHED_SEQUENCE_LENGTHS' (%d), was '%d'",
				sequence.MISSMATCHED_SEQUENCE_LENGTHS, is.Errno)
		}
	default:
		t.Errorf("Error type was wrong %g", is)
	}
}

func TestTwoSpiecesWithSpecialCharacters(t *testing.T) {
	sequence1 := sequence.NewSequence("Homo sapiens", []byte("ATAGCT[AC]G"))
	sequence1.Species = "Homo sapiens"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("Homo erectus", []byte("ATAGCTAC"))
	sequence2.Species = "Homo erectus"
	sequence2.Gene = "ATP8"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2)

	buf := bytes.Buffer{}

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Error("Expected no error, got one", err)
	}

	tnt.WriteXRead(&buf)

	expected := `xread
'Title Here'
8 2
Homo_erectus ATAGCTAC
Homo_sapiens ATAGCT[AC]G
;`
	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}
}

// TestTwoGenesTwoSpecies should sort the blocks by alphabetical order
func TestTwoGenesTwoSpecies(t *testing.T) {
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
;`

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2, sequence3, sequence4)

	buf := bytes.Buffer{}

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Error("Expected no error, got one", err)
	}

	if err := tnt.WriteXRead(&buf); err != nil {
		t.Error("Expected no error, got one", err)
	}

	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}

}

func TestFourSpiecesWithDefinedOutgroupFromMiddle(t *testing.T) {
	sequence1 := sequence.NewSequence("A a", []byte("ATAGCTACG"))
	sequence1.Species = "A a"
	sequence1.Gene = "ATP8"

	sequence2 := sequence.NewSequence("B b", []byte("ATAGTCACG"))
	sequence2.Species = "B b"
	sequence2.Gene = "ATP8"

	sequence3 := sequence.NewSequence("C c", []byte("ATAGCTACG"))
	sequence3.Species = "C c"
	sequence3.Gene = "ATP8"

	sequence4 := sequence.NewSequence("D d", []byte("ATAGTCACG"))
	sequence4.Species = "D d"
	sequence4.Gene = "ATP8"

	tnt := &formats.TNT{Title: "Title Here"}
	tnt.AddSequence(sequence1, sequence2, sequence3, sequence4)
	tnt.SetOutgroup("C c")

	buf := bytes.Buffer{}

	if _, err := tnt.GenerateMetaData(); err != nil {
		t.Error("Expected no error, got one", err)
	}

	tnt.WriteXRead(&buf)

	expected := `xread
'Title Here'
9 4
C_c ATAGCTACG
A_a ATAGCTACG
B_b ATAGTCACG
D_d ATAGTCACG
;`
	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n\n\"%s\"\n\nGot:\n\n\"%s\"", expected, got)
	}
}
