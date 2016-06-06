package sequence

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"

	"github.com/yarbelk/refasta/scanner"
)

// Sequence is the base nucleotide container, it has a Name
// which is the full descriptive name of the sequence
// and Seq (byte), which is the nucleotide data
type Sequence struct {
	// Generic Identifier
	Name string
	// Species Name
	Species string
	// Gene name
	Gene string
	Seq  SequenceData
	// Length is the Logical Length of a sequence; where polymorphics are counted
	// as a single length eg: [AG] == len 1
	Length int
	// alphabet is the unique character in this sequence
	alphabet map[rune]bool
}

var safeRegex = regexp.MustCompile("( )")

func Safe(in string) string {
	return safeRegex.ReplaceAllLiteralString(in, "_")
}

// seq is a specific type of byte slice, so we can throw
// methods on there for easy handling
type SequenceData []byte

// NewSequence returns a value type Sequence, this will scan the sequence data
// to determin its length; so only use this is you haven't already done something
// to get the length of the data.
func NewSequence(name string, seq []byte) Sequence {
	_, length, alpha, _ := scanner.ScanSequenceData(bufio.NewReader(bytes.NewBuffer(seq)))
	return Sequence{Name: name, Seq: SequenceData(seq), Length: length, alphabet: alpha}
}

// String represesntiation of a Seq is just typecasting to a `string`
func (s SequenceData) String() string {
	return string(s)
}

// SafeName will replace spaces with underscores (possibly other things in
// the future as I find the need
func (s Sequence) SafeName() string {
	return Safe(s.Name)
}

// SafeSpecies repalces spaces in the species name with underscores
func (s Sequence) SafeSpecies() string {
	return Safe(s.Species)
}

// Repr returns a representation of the data
func (s Sequence) GoString() string {
	var truncatedSequence []byte
	if len(s.Seq) > 5 {
		truncatedSequence = append(s.Seq[:5], []byte("...")...)
	} else {
		truncatedSequence = s.Seq[:]
	}
	return fmt.Sprintf(
		"Sequence{Name: %s, Species %s. Seqence: [%s], Length: %d}",
		s.Name,
		s.Species,
		truncatedSequence,
		s.Length,
	)
}

// SetAlphabet will set the alphabet map.
func (s *Sequence) SetAlphabet(alpha map[rune]bool) {
	s.alphabet = alpha
}
