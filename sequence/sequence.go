package sequence

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"github.com/yarbelk/refasta/scanner"
)

type ErrNo int

const (
	UNKNOWN ErrNo = iota
	MISSMATCHED_SEQUENCE_LENGTHS
)

// Sequence is the base nucleotide container, it has a Name
// which is the full descriptive name of the sequence
// and Seq (byte), which is the nucleotide data
type Sequence struct {
	Name    string
	Species string
	Gene    string
	Seq     SequenceData
	Length  int
}

// GeneMetaData is a convience structure, which is useful for sorting and keying
// and validating data arround sequences and genes.  Mostly useful because
// it is a really small structure, and we can implement a bunch of methods on it
// so we don't have to worry about the memory footprint of giant byte arrays.
type GeneMetaData struct {
	Gene          string
	Length        int
	NumberSpecies int
}

// GMDSlice is a GeneMetaData slice, which implements the sort.Interface
type GMDSlice []GeneMetaData

// Len returns the length of the underlying slice, part of the sort.Interface
func (g GMDSlice) Len() int {
	return len(g)
}

// Less compairs the Gene name, part of the sort.Interface
func (g GMDSlice) Less(i, j int) bool {
	return g[i].Gene < g[j].Gene
}

// Sort is a convienence method for sorting a GMDSlice
func (g GMDSlice) Sort() {
	sort.Sort(g)
}

// Swap swaps the positions of i, j; part of the sort.Interface
func (g GMDSlice) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

var safeRegex = regexp.MustCompile("( )")

type InvalidSequence struct {
	Message string
	Details string
	Errno   ErrNo
}

func (e InvalidSequence) Error() string {
	return fmt.Sprintf("InvalidSequence: %s\nDetails: %s")
}

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
	_, length, _ := scanner.ScanSequenceData(bufio.NewReader(bytes.NewBuffer(seq)))
	return Sequence{Name: name, Seq: SequenceData(seq), Length: length}
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
