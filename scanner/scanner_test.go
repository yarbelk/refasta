package scanner_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/yarbelk/refasta/scanner"
)

func TestStoresLengthOfNormalSequence(t *testing.T) {
	buf := bufio.NewReader(bytes.NewBuffer([]byte("ATAGC")))
	lit, length, alpha, err := scanner.ScanSequenceData(buf)
	expected := 5
	if err != nil {
		t.Errorf("Expected no error: got '%s'", err.Error())
	}
	if len(alpha) != 4 {
		t.Errorf("Alphabet returned didn't have the right length; expected 4, got %d", len(alpha))
	}
	if string(lit) != "ATAGC" {
		t.Errorf("Expected: '%s', got '%s'", "ATAGC", lit)
	}
	if length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, length)
	}
}

func TestStoresLengthOfSequenceWithDash(t *testing.T) {
	buf := bufio.NewReader(bytes.NewBuffer([]byte("AT-AGC")))
	lit, length, _, err := scanner.ScanSequenceData(buf)
	expected := 6
	if err != nil {
		t.Errorf("Expected no error: got '%s'", err.Error())
	}
	if string(lit) != "AT-AGC" {
		t.Errorf("Expected: '%s', got '%s'", "AT-AGC", lit)
	}
	if length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, length)
	}
}

func TestStoresLengthOfSequenceWithQuestionMark(t *testing.T) {
	buf := bufio.NewReader(bytes.NewBuffer([]byte("AT?AGC")))
	_, length, _, err := scanner.ScanSequenceData(buf)
	expected := 6
	if err != nil {
		t.Errorf("Expected no error: got '%s'", err.Error())
	}
	if length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, length)
	}
}

func TestStoresLengthOfSequenceWithBraces(t *testing.T) {
	buf := bufio.NewReader(bytes.NewBuffer([]byte("AT[AG]C")))
	_, length, _, err := scanner.ScanSequenceData(buf)
	expected := 4
	if err != nil {
		t.Errorf("Expected no error: got '%s'", err.Error())
	}
	if length != expected {
		t.Errorf("Expected: '%d', got '%d'", expected, length)
	}
}
