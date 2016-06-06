package sequence

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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

// SequenceType is one of used in TNT (and probably other things)
// to determine if the sequence is a DNA sequence or protein sequence
type SequenceType int

const (
	UNSUPPORTED_TYPE SequenceType = iota
	DNA_TYPE
	PROTEIN_TYPE
	BLANK_TYPE
)

const (
	PROTEIN_ALPHABET = "GALMFWKQESPVICYHRNDT"
	DNA_ALPHABET     = "ACGT"
)

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

// isProtein will return true if this is a protein
// it will log a warning if there are a small number of amino acids
// I don't particularly like making this a method, but its the easiest
// way to output a warning
//
// Returns TRUE if positivily a protein
func (s *Sequence) isProtein(alphabet map[rune]bool) bool {
	var c int
	for _, a := range PROTEIN_ALPHABET {
		switch {
		case alphabet[a]:
			c++
		case a == 'X':
			c++
		case a == '-':
		case a == '?':
		default:
			continue
		}
	}

	protein := c > 5 // completly arbitrary: 5 > len(DNA_ALPHABET)
	if protein && c < len(PROTEIN_ALPHABET) {
		fmt.Fprintf(os.Stderr, "%s only has %d Amino Acids, please check your data set\n", s.GoString(), c)
	}
	return protein
}

// isBlank checks to see if its a blank type
func (s *Sequence) isBlank(alphabet map[rune]bool) bool {
	for _, c := range s.Seq {
		if c != '-' {
			return false
		}
	}
	return true
}

// isDNA checks to see if a sequence is DNA
func (s *Sequence) isDNA(alphabet map[rune]bool) bool {
	var c int
	if s.isProtein(alphabet) {
		return false
	}
	for _, a := range DNA_ALPHABET {
		if alphabet[a] {
			c++
		}
	}
	dna := c >= 3 // Arbitrary
	if dna && c < len(DNA_ALPHABET) {
		fmt.Fprintf(os.Stderr, "%s only has %d Nucleic Acids, please check your data set\n", s.GoString(), c)
	}
	return dna
}

// Type of sequence, DNA, Protein, or Unsupported
// TODO: This should be pulled into a single loop!
func (s *Sequence) Type() SequenceType {
	switch {
	case s.isBlank(s.alphabet):
		return BLANK_TYPE
	case s.isProtein(s.alphabet):
		return PROTEIN_TYPE
	case s.isDNA(s.alphabet):
		return DNA_TYPE
	default:
		fmt.Fprintf(os.Stderr, "Couldn't determine sequence type, %s\n", s.GoString())
		return UNSUPPORTED_TYPE
	}
}
