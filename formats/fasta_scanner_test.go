package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
)

func TestFastaScanReturnsSequenceName(t *testing.T) {
	reader := bytes.NewReader([]byte(">Sequence Identifier\n"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, _, length := fastaScanner.Scan()

	if tok != formats.SEQUENCE_ID {
		t.Errorf("token should have been 'formats.SEQUENCE_ID' %d, was '%d'", formats.SEQUENCE_ID, tok)
	}

	expectedLit := "Sequence Identifier"
	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 19 {
		t.Error("Expected sequence length to be '19', got '%d'", length)
	}
}

func TestInvalidCharacterReturnsScanTokenInvalid(t *testing.T) {
	reader := bytes.NewReader([]byte("𠜎"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, _, length := fastaScanner.Scan()

	if tok != formats.INVALID {
		t.Errorf("token should have been 'formats.INVALID' %d, was '%d'", formats.INVALID, tok)
	}

	expectedLit := ""
	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 0 {
		t.Error("Expected sequence length to be '0', got '%d'", length)
	}
}

func TestFastaScanScansNucleotideData(t *testing.T) {
	reader := bytes.NewReader([]byte("ATGCGTA"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, alpha, length := fastaScanner.Scan()

	if len(alpha) != 4 {
		t.Errorf("Wrong alphabet length, should be 4, was %d", len(alpha))
	}
	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 7 {
		t.Error("Expected sequence length to be '7', got '%d'", length)
	}
}

func TestFastaScanScansNucleotideDataWithDash(t *testing.T) {
	expectedLit := "ATG-CGTA"
	reader := bytes.NewReader([]byte(expectedLit))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, _, length := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 8 {
		t.Error("Expected sequence length to be '8', got '%d'", length)
	}
}

func TestFastaScanScansNucleotideDataWithNewline(t *testing.T) {
	reader := bytes.NewReader([]byte("ATG\nCGTA"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, _, length := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 7 {
		t.Error("Expected sequence length to be '7', got '%d'", length)
	}
}

func TestFastaScanScansNucleotideDataWithNewlineAndNewSequence(t *testing.T) {
	reader := bytes.NewReader([]byte("ATG\nCGTA\n>Bob"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit, _, length := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if string(lit) != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}

	if length != 7 {
		t.Error("Expected sequence length to be '7', got '%d'", length)
	}
}

func TestFastaScanTwoSequencesInTotal(t *testing.T) {
	type expected struct {
		Tok    formats.Token
		Data   string
		Length int
	}
	expecteds := []expected{
		expected{formats.SEQUENCE_ID, "foo", 3},
		expected{formats.SEQUENCE_DATA, "ATGCGTA", 7},
		expected{formats.SEQUENCE_ID, "bar", 3},
		expected{formats.SEQUENCE_DATA, "TATGCGTAT", 9},
	}
	reader := bytes.NewReader([]byte(">foo\nATG\nCGTA\n>bar\nTATGC\nGTAT"))
	fastaScanner := formats.NewFastaScanner(reader)

	for _, v := range expecteds {
		tok, lit, _, length := fastaScanner.Scan()
		if tok != v.Tok {
			t.Errorf("token should have been '%d', was '%d'", v.Tok, tok)
		}

		if string(lit) != v.Data {
			t.Errorf("value should be '%s', was '%s'", v.Data, string(lit))
		}

		if length != v.Length {
			t.Errorf("Expected sequence length to be '%d', got '%d'", v.Length, length)
		}
	}
}
