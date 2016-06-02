package formats_test

import (
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/formats"
)

func TestFastaScanReturnsSequenceName(t *testing.T) {
	reader := bytes.NewReader([]byte(">Sequence Identifier\n"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit := fastaScanner.Scan()

	if tok != formats.SEQUENCE_ID {
		t.Errorf("token should have been 'formats.SEQUENCE_ID' %d, was '%d'", formats.SEQUENCE_ID, tok)
	}

	expectedLit := "Sequence Identifier"
	if lit != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}
}

func TestInvalidCharacterReturnsScanTokenInvalid(t *testing.T) {
	reader := bytes.NewReader([]byte("𠜎"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit := fastaScanner.Scan()

	if tok != formats.INVALID {
		t.Errorf("token should have been 'formats.INVALID' %d, was '%d'", formats.INVALID, tok)
	}

	expectedLit := ""
	if lit != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}
}

func TestFastaScanScansNucleotideData(t *testing.T) {
	reader := bytes.NewReader([]byte("ATGCGTA"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if lit != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}
}

func TestFastaScanScansNucleotideDataWithNewline(t *testing.T) {
	reader := bytes.NewReader([]byte("ATG\nCGTA"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if lit != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}
}

func TestFastaScanScansNucleotideDataWithNewlineAndNewSequence(t *testing.T) {
	reader := bytes.NewReader([]byte("ATG\nCGTA\n>Bob"))
	fastaScanner := formats.NewFastaScanner(reader)
	tok, lit := fastaScanner.Scan()

	if tok != formats.SEQUENCE_DATA {
		t.Errorf("token should have been 'formats.SEQUENCE_DATA' %d, was '%d'", formats.SEQUENCE_DATA, tok)
	}

	expectedLit := "ATGCGTA"
	if lit != expectedLit {
		t.Errorf("Sequence should be '%s', was '%s'", expectedLit, lit)
	}
}
