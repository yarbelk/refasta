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
	fastaFormat      = ">%s\n%s\n"
	testSequenceName = "Oryg luct 1033"
)

var testSequence = []byte("TACTTTAGCGGTAATGTATAGGTATACAACTTAAGCGCCATCCATTTTAAGGGCTAGTTGCTTCGGCAGGTGAGTTGTTACACACTCCTTAGCGGATTACGACTTCCATGTCCACCGTCCTGCTGTTTTAAGCAACCAACGCCTTTCATGGTATCTGCATGAGTTGTTAATTTGGGCACCGTAACATTACGTTTGGTTCATCCCACAGCGCCAGTTCTGCTTACCAAAAGTGGCCCACTGGGCACATTATATCATAACCACAACCTTCATATCAAGTAAGGTGAGGTTCTTACCCATTTAAAGTTTGAGA")

// given a sequence, we should be able to write it to a file
// Assume no interleaving at this point
func TestWithSequenceCanWriteOneFasta(t *testing.T) {
	sequence := sequence.NewSequence(testSequenceName, testSequence)
	expected := fmt.Sprintf(fastaFormat, "Oryg_luct_1033", testSequence) + "\n"

	output := &bytes.Buffer{}

	fastaWriter := formats.FastaWriter{
		FileName: "test_out.fasta",
		File:     output}
	fastaWriter.AddSequence(sequence)
	fastaWriter.WriteSequences()
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

	fastaWriter := formats.FastaWriter{
		FileName: "test_out.fasta",
		File:     output}
	fastaWriter.AddSequence(sequence1)
	fastaWriter.AddSequence(sequence2)
	fastaWriter.WriteSequences()
	got := output.String()

	if got != expected {
		t.Errorf("Did not get expected outputs for basic write fasta:\n\n\tGot:\n\n%s\n\n\tExpected:\n\n%s", got, expected)
	}
}
