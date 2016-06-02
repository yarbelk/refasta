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
