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

func TestStoresLengthOfNormalSequence(t *testing.T) {
	seq := sequence.NewSequence("name with spaces", []byte("ATAGC"))
	expected := 5
	if seq.Length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, seq.Length)
	}
}

func TestStoresLengthOfSequenceWithDash(t *testing.T) {
	seq := sequence.NewSequence("name with spaces", []byte("AT-AGC"))
	expected := 6
	if seq.Length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, seq.Length)
	}
}

func TestStoresLengthOfSequenceWithQuestionMark(t *testing.T) {
	seq := sequence.NewSequence("name with spaces", []byte("AT?AGC"))
	expected := 6
	if seq.Length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, seq.Length)
	}
}

func TestStoresLengthOfSequenceWithBraces(t *testing.T) {
	seq := sequence.NewSequence("name with spaces", []byte("AT[AG]C"))
	expected := 4
	if seq.Length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, seq.Length)
	}
}

func TestSequenceLengthSetting(t *testing.T) {
	seq := sequence.NewSequence("test", []byte("asdf"))
	if seq.Length != 0 {
		t.Errorf("Initialization length of sequence should be zero")
	}
	seq.Length = 1
	if seq.Length != 1 {
		t.Errorf("seq should now have a Length of '1', was '%d'", seq.Length)
	}
}
