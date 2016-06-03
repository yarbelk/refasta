package sequence

import "sort"

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
