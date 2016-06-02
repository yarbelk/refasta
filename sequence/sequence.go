package sequence

import "regexp"

// Sequence is the base nucleotide container, it has a Name
// which is the full descriptive name of the sequence
// and Seq (byte), which is the nucleotide data
type Sequence struct {
	Name    string
	Species string
	Gene    string
	Seq     sequenceData
}

var safeRegex = regexp.MustCompile("( )")

func Safe(in string) string {
	return safeRegex.ReplaceAllLiteralString(in, "_")
}

// seq is a specific type of byte slice, so we can throw
// methods on there for easy handling
type sequenceData []byte

// NewSequence returns a value type Sequence
func NewSequence(name string, seq []byte) Sequence {
	return Sequence{Name: name, Seq: sequenceData(seq)}
}

// String represesntiation of a Seq is just typecasting to a `string`
func (s sequenceData) String() string {
	return string(s)
}

// SafeName will replace spaces with underscores (possibly other things in
// the future as I find the need
func (s Sequence) SafeName() string {
	return Safe(s.Name)
}

func (s Sequence) SafeSpecies() string {
	return Safe(s.Species)
}
