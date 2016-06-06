package formats_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/yarbelk/refasta/formats"
	"github.com/yarbelk/refasta/sequence"
)

/* here is a DNA sequence!
>Oryg_luct_1033
TACTTTAGCGGTAATGTATAGGTATACAACTTAAGCGCCATCCATTTTAAGGGCTAGTTGCTTCGGCAGGTGAGTTGTT
ACACACTCCTTAGCGGATTACGACTTCCATGTCCACCGTCCTGCTGTTTTAAGCAACCAACGCCTTTCATGGTATCTGC
ATGAGTTGTTAATTTGGGCACCGTAACATTACGTTTGGTTCATCCCACAGCGCCAGTTCTGCTTACCAAAAGTGGCCCA
CTGGGCACATTATATCATAACCACAACCTTCATATCAAGTAAGGTGAGGTTCTTACCCATTTAAAGTTTGAGA
*/

const (
	testGeneName     = "testGeneName"
	fastaFormat      = ">%s\n%s\n"
	testSequenceName = "Oryg_luct_1033"
)

var testSequence = []byte("TACTTTAGCGGTAATGTATAGGTATACAACTTAAGCGCCATCCATTTTAAGGGCTAGTTGCTTCGGCAGGTGAGTTGTTACACACTCCTTAGCGGATTACGACTTCCATGTCCACCGTCCTGCTGTTTTAAGCAACCAACGCCTTTCATGGTATCTGCATGAGTTGTTAATTTGGGCACCGTAACATTACGTTTGGTTCATCCCACAGCGCCAGTTCTGCTTACCAAAAGTGGCCCACTGGGCACATTATATCATAACCACAACCTTCATATCAAGTAAGGTGAGGTTCTTACCCATTTAAAGTTTGAGA")

// given a sequence, we should be able to write it to a file
// Assume no interleaving at this point
func TestWithSequenceCanWriteOneFasta(t *testing.T) {
	sequence := sequence.NewSequence(testSequenceName, testSequence)
	expected := fmt.Sprintf(fastaFormat, "Oryg_luct_1033", testSequence) + "\n"

	output := &bytes.Buffer{}

	fasta := formats.Fasta{}
	fasta.AddSequence(sequence)
	fasta.WriteSequences(output)
	got := output.String()
	if got != expected {
		t.Errorf("Did not get expected outputs for basic write fasta:\n\n\tGot:\n\n%s\n\n\tExpected:\n\n%s", got, expected)
	}
}

func TestCanWriteTwoFasta(t *testing.T) {
	sequence1 := sequence.NewSequence(testSequenceName, testSequence)
	sequence2 := sequence.NewSequence("Sequence Two", testSequence)
	expected := fmt.Sprintf(fastaFormat, "Oryg_luct_1033", testSequence) +
		fmt.Sprintf(fastaFormat, "Sequence_Two", testSequence) + "\n"

	output := &bytes.Buffer{}

	fastaWriter := formats.Fasta{}
	fastaWriter.AddSequence(sequence1)
	fastaWriter.AddSequence(sequence2)
	fastaWriter.WriteSequences(output)
	got := output.String()

	if got != expected {
		t.Errorf("Did not get expected outputs for basic write fasta:\n\n\tGot:\n\n%s\n\n\tExpected:\n\n%s", got, expected)
	}
}

func TestCanParseNonInterleavedSingleSequence(t *testing.T) {
	inputString := fmt.Sprintf(fastaFormat, "Oryg_luct_1033", testSequence) + "\n"
	input := bytes.NewBuffer([]byte(inputString))
	fastaReader := formats.Fasta{}
	fastaReader.Parse(input)
	expectedParsed := 1
	if len(fastaReader.Sequences) != expectedParsed {
		t.Errorf("Wrong number of sequences after parseing; expected: '%d', got '%d'", expectedParsed, len(fastaReader.Sequences))
	}

	seq := fastaReader.Sequences[0]
	expectedName := "Oryg_luct_1033"

	if seq.Name != expectedName {
		t.Errorf("Expected name didn't match: expected '%s', got '%s'", expectedName, seq.Name)
	}

	if string(seq.Seq) != string(testSequence) {
		t.Errorf("Expected name didn't match: expected '%s', got '%s'", string(testSequence), string(seq.Seq))
	}
}

func TestGetGeneNameFromFileName(t *testing.T) {
	inputString := fmt.Sprintf(fastaFormat, "Oryg_luct_1033", testSequence) + "\n"
	input := bytes.NewBuffer([]byte(inputString))
	fastaReader := formats.Fasta{}

	fastaReader.Parse(input, testGeneName)
	expectedParsed := 1

	if len(fastaReader.Sequences) != expectedParsed {
		t.Errorf("Wrong number of sequences after parseing; expected: '%d', got '%d'", expectedParsed, len(fastaReader.Sequences))
	}

	seq := fastaReader.Sequences[0]
	if seq.Gene != testGeneName {
		t.Errorf("Default gene name for fasta didn't match '%s' name, was '%s'", testGeneName, seq.Gene)
	}
}

func TestDataBeforeIDReturnsError(t *testing.T) {
	inputString := "ATAG\n>shouldFail"
	input := bytes.NewBuffer([]byte(inputString))
	fastaReader := formats.Fasta{}

	err := fastaReader.Parse(input, testGeneName)

	if err == nil {
		t.Errorf("Expected error to be; error was <nil>")
	}

	if err.(sequence.FormatError).Errno != sequence.BAD_FORMAT {
		t.Errorf("Expected Errno to be '%d', was '%d'",
			sequence.BAD_FORMAT,
			err.(sequence.FormatError).Errno)
	}
}

func TestTwoIDSBackToBackReturnsError(t *testing.T) {
	inputString := ">IDOne\n>shouldFail"
	input := bytes.NewBuffer([]byte(inputString))
	fastaReader := formats.Fasta{}

	err := fastaReader.Parse(input, testGeneName)

	if err == nil {
		t.Errorf("Expected error to be; error was <nil>")
	}

	if err.(sequence.FormatError).Errno != sequence.BAD_FORMAT {
		t.Errorf("Expected Errno to be '%d', was '%d'",
			sequence.BAD_FORMAT,
			err.(sequence.FormatError).Errno)
	}
}
