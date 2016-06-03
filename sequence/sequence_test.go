package sequence_test

import (
	"testing"

	"github.com/yarbelk/refasta/sequence"
)

func TestSafeName(t *testing.T) {
	seq := sequence.NewSequence("name with spaces", []byte("asdf"))
	expected := "name_with_spaces"
	if seq.SafeName() != expected {
		t.Errorf("Expected: '%q', got '%q'", expected, seq.SafeName())
	}
}

func TestSafeSpace(t *testing.T) {
	in := "Hello World!"
	expected := "Hello_World!"
	got := sequence.Safe(in)

	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
	}
}

func TestNewSequenceCountsLength(t *testing.T) {
	seq := sequence.NewSequence("test", []byte("ATAGATAG"))
	expected := 8
	if seq.Length != expected {
		t.Errorf("sequence has wrong length, expected '%d', got '%d'", expected, seq.Length)
	}
}

func TestNewSequenceCountsLengthWithBrackets(t *testing.T) {
	seq := sequence.NewSequence("test", []byte("ATAGAT[AG]"))
	expected := 7
	if seq.Length != expected {
		t.Errorf("sequence has wrong length, expected '%d', got '%d'", expected, seq.Length)
	}
}
