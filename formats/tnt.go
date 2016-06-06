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
	Title        string
	Sequences    map[string]map[string]sequence.Sequence
	MetaData     sequence.GMDSlice
	speciesNames []string
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

// Construct a species using a GMDSlice to order the gene sequences
func (t *TNT) PrintableTaxa() []taxonData {
	t.MetaData.Sort()
	var allSpecies []taxonData = make([]taxonData, 0, len(t.speciesNames))

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
	return allSpecies
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
	allSpecies := t.PrintableTaxa()
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

	if err := t.WriteXRead(writer); err != nil {
		return err
	}

	if err := t.WriteBlocks(writer); err != nil {
		return err
	}

	return nil
}

// GenerateMetaData will make sure that the sequences for the same
// gene sequence (or whatever sequence) are all the same length.
// Returns types of InvalidSequence with ErrNo
// MISSMATCHED_SEQUENCE_LENGTHS if they are no correct
// If they are correct, it will return a slice of the gene meta data
// GeneMetaData, sequence.GMDSlice
func (t *TNT) GenerateMetaData() (sequence.GMDSlice, error) {
	var expectedLen int
	geneMetaData := make(sequence.GMDSlice, 0, len(t.Sequences))

	for gene, _ := range t.Sequences {
		for i, name := range t.speciesNames {
			seq := t.Sequences[gene][name]
			if i == 0 {
				expectedLen = seq.Length
			} else if seq.Length != expectedLen {
				return nil, sequence.InvalidSequence{
					Message: "Sequences are not the Same length",
					Details: "None so far",
					Errno:   sequence.MISSMATCHED_SEQUENCE_LENGTHS,
				}
			}
		}
		geneMetaData = append(geneMetaData, sequence.GeneMetaData{
			Gene:          gene,
			Length:        expectedLen,
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
