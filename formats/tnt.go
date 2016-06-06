package formats

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/template"
	"github.com/yarbelk/refasta/sequence"
)

// TNT formatter
type TNT struct {
	Title             string
	Sequences         map[string]map[string]sequence.Sequence
	MetaData          sequence.GMDSlice
	speciesNames      []string
	Outgroup          string
	dirtyData         bool
	maxSequenceLength int
	blankSeq          sequence.SequenceData
}

const tntNonInterleavedTemplateString = `xread
'{{ .Title }}'
{{ .Length }} {{ .NTaxa }}
{{ range $i, $taxon := .Taxa }}{{ $taxon.SpeciesName }} {{ $taxon.Sequence }}
{{ end }};`

const tntBlocksTemplateString = `
blocks {{ .Blocks }};
cnames
{{ range $i, $cname := .Cnames}}{{ $cname }}
{{ end }};`

var tntNonInterleavedTemplate = template.Must(template.New("TNTXread").Parse(tntNonInterleavedTemplateString))
var tntBlocksTemplate = template.Must(template.New("TNTBlocks").Parse(tntBlocksTemplateString))

type templateContext struct {
	Title         string
	Length, NTaxa int
	Taxa          []taxonData
}

type taxonData struct {
	SpeciesName string
	Sequence    sequence.SequenceData
}

const TNT_FORMAT = "tnt"

/*
Construct a species using a GMDSlice to order the gene sequences.
If there is a defined outgroup, then sort that to the front of the
printable list
*/
func (t *TNT) PrintableTaxa() ([]taxonData, error) {
	if t.MetaData == nil {
		if _, err := t.GenerateMetaData(); err != nil {
			return nil, err
		}
	}
	var allSpecies []taxonData = make([]taxonData, 0, len(t.speciesNames))
	t.sortByOutgroup()

	for _, n := range t.speciesNames {
		combinedSequences := make([]byte, 0, t.getTotalLength())
		for _, gmd := range t.MetaData {
			combinedSequences = append(combinedSequences, t.Sequences[gmd.Gene][n].Seq...)
		}
		allSpecies = append(allSpecies, taxonData{
			SpeciesName: sequence.Safe(n),
			Sequence:    combinedSequences,
		})
	}
	return allSpecies, nil
}

/*
sortByOutgroup is a helper method to sort the species names,
starting with the outgroup.  This is used to format the xread block
with the outgroup as the first species in the list.
*/
func (t *TNT) sortByOutgroup() {
	if t.Outgroup == "" {
		return
	}
	var index int = 0
	safeOG := sequence.Safe(t.Outgroup)
	for j, n := range t.speciesNames {
		if safeOG == sequence.Safe(n) {
			index = j
			break
		}
	}

	a := t.speciesNames[index]
	t.speciesNames = append(
		append(t.speciesNames[:1], t.speciesNames[0:index]...),
		t.speciesNames[index+1:]...)
	t.speciesNames[0] = a
}

// insertString into the place that would keep it uniquely and ordered ascending
func insertString(slice []string, s string) []string {
	i := sort.SearchStrings(slice, s)
	// Inserstion sort of the species names: builds up the list as a sorted list
	if i < len(slice) && slice[i] != s {
		// Species Name not in the list; insert it at i
		slice = append(slice[:i], append([]string{s}, slice[i:]...)...)
	} else if i == len(slice) {
		slice = append(slice, s)
	}
	return slice
}

// AddSequence (or multiple) to the internal sequence store.
func (t *TNT) AddSequence(seqs ...sequence.Sequence) {
	for _, seq := range seqs {
		if t.Sequences == nil {
			t.Sequences = make(map[string]map[string]sequence.Sequence)
		}
		if m, ok := t.Sequences[seq.Gene]; !ok || m == nil {
			t.Sequences[seq.Gene] = make(map[string]sequence.Sequence)
		}
		t.Sequences[seq.Gene][seq.Species] = seq
		t.speciesNames = insertString(t.speciesNames, seq.Species)
	}
}

/*
WriteXRead writes out the xread block; which contains the sequence
and taxa data

	xread
	taxa_1 CTAGC...
	taxa_2 TAGCA...
	;

*/
func (t *TNT) WriteXRead(writer io.Writer) error {
	allSpecies, err := t.PrintableTaxa()
	if err != nil {
		return err
	}
	context := templateContext{
		Title:  t.Title,
		Length: t.getTotalLength(),
		NTaxa:  len(t.speciesNames),
		Taxa:   allSpecies,
	}
	return tntNonInterleavedTemplate.Execute(writer, context)
}

/*
WriteBlocks writes out the block definitions and their names.

	blocks 0 10 18 200;
	cnames
	[1 ATP8;
	[2 ATP6;
	[3 co1;
	[4 dblsex;
	;

There is an implicit block `[0 "ALL"`, which cannot be renamed,
so the first user defined block is `1`.
*/

func (t *TNT) WriteBlocks(writer io.Writer) error {
	var startPos []string = make([]string, 0, len(t.MetaData))
	var cnames []string = make([]string, 0, len(t.MetaData))

	var newStart int
	for i, _ := range t.MetaData {
		if i != 0 {
			newStart = newStart + t.MetaData[i-1].Length
		}
		cname := fmt.Sprintf("[%d %s;", i+1, t.MetaData[i].Gene)

		cnames = append(cnames, cname)
		startPos = append(startPos, strconv.Itoa(newStart))
	}
	blocks := strings.Join(startPos, " ")
	context := struct {
		Blocks string
		Cnames []string
	}{
		Blocks: blocks,
		Cnames: cnames,
	}
	return tntBlocksTemplate.Execute(writer, context)
}

// WriteSequences will collect up the sequences, verify their validity,
// and output a formated TNT file to the supplied writer
func (t *TNT) WriteSequences(writer io.Writer) error {
	if _, err := t.GenerateMetaData(); err != nil {
		return err
	}
	if t.dirtyData {
		t.CleanData()
	}

	if err := t.WriteXRead(writer); err != nil {
		return err
	}

	if err := t.WriteBlocks(writer); err != nil {
		return err
	}

	return nil
}

func geneLength(lengths map[int][]string) (max int) {
	for i, _ := range lengths {
		if i > max {
			max = i
		}
	}
	return
}

// fmtInvalidSequenceErr will return a specialized error for invalid
// sequence lengths.
func fmtInvalidSequenceErr(sequenceName string, lengths map[int][]string) error {
	details := []string{}
	for length, seqs := range lengths {
		details = append(details, fmt.Sprintf("\t%d: %s", length, strings.Join(seqs, ", ")))
	}

	detailedMessage := fmt.Sprintf("Sequence %s has inconsistant sequence lengths:\n%s", sequenceName, strings.Join(details, "\n"))
	return sequence.InvalidSequence{
		Message: "Sequences are not the Same length",
		Details: detailedMessage,
		Errno:   sequence.MISSMATCHED_SEQUENCE_LENGTHS,
	}
}

/*
GenerateMetaData will make sure that the sequences for the same
gene sequence (or whatever sequence) are all the same length.
Returns types of InvalidSequence with ErrNo
MISSMATCHED_SEQUENCE_LENGTHS if they are no correct
If they are correct, it will return a slice of the gene meta data
GeneMetaData, sequence.GMDSlice

If a sequence is zero; it is not counted as bad.  It  needs to be
cleaned up with a call to CleanData

This will also set the max lenght sequence size; which is used
by some helper functions
*/
func (t *TNT) GenerateMetaData() (sequence.GMDSlice, error) {
	geneMetaData := make(sequence.GMDSlice, 0, len(t.Sequences))

	for gene, _ := range t.Sequences {
		lengths := make(map[int][]string)
		for _, name := range t.speciesNames {
			seq := t.Sequences[gene][name]
			if seq.Length > t.maxSequenceLength {
				t.maxSequenceLength = seq.Length
			}
			if _, ok := lengths[seq.Length]; ok {
				lengths[seq.Length] = append(lengths[seq.Length], seq.Name)
			} else {
				lengths[seq.Length] = []string{seq.Name}
			}
		}
		_, hasZero := lengths[0]
		if (len(lengths) > 2) || (len(lengths) > 1 && !hasZero) {
			return nil, fmtInvalidSequenceErr(gene, lengths)
		}

		if hasZero {
			t.dirtyData = true
		}
		geneMetaData = append(geneMetaData, sequence.GeneMetaData{
			Gene:          gene,
			Length:        geneLength(lengths),
			NumberSpecies: len(t.Sequences[gene]),
		})
	}
	t.MetaData = geneMetaData
	geneMetaData.Sort()
	return geneMetaData, nil
}

// getTotalLength will return the combined length of all genes.  This should
// be the same for each species.
func (t *TNT) getTotalLength() (length int) {
	for _, gmd := range t.MetaData {
		length = length + gmd.Length
	}
	return
}

// SetOutgroup will set the outgroup under test.  This sorts it to the top
// of the list of taxa in the xread block
func (t *TNT) SetOutgroup(species string) error {
	t.Outgroup = species
	return nil
}

/*
blankSequence returns a slice of '---' bytes.  Does this by pre-allocating

the longest that could be returned, and slicing up subsets of it to be returned.
this should be fine because once returned, they should never be modified,
so the shared memory should not be a problem
*/
func (t *TNT) blankSequence(n int) sequence.SequenceData {
	if t.blankSeq == nil || len(t.blankSeq) < t.maxSequenceLength {
		t.blankSeq = make(sequence.SequenceData, t.maxSequenceLength, t.maxSequenceLength)
		for i, _ := range t.blankSeq {
			t.blankSeq[i] = '-'
		}
	}

	return t.blankSeq[:n]
}

/*
CleanData will fill in missing data.
*/
func (t *TNT) CleanData() {
	for _, gmd := range t.MetaData {
		for _, name := range t.speciesNames {
			seq := t.Sequences[gmd.Gene][name]
			if len(seq.Seq) == 0 {
				seq.Seq = t.blankSequence(gmd.Length)
				seq.Length = gmd.Length
				t.Sequences[gmd.Gene][name] = seq
			}
		}
	}

}
