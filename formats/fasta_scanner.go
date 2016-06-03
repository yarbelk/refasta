package formats

import (
	"bufio"
	"io"

	"github.com/yarbelk/refasta/scanner"
)

const (
	UNSTARTED Token = -1
	EOF       Token = iota
	SEQUENCE_ID
	SEQUENCE_DATA
	WHITESPACE
	INVALID
)

type Token int

type FastaScanner struct {
	reader *bufio.Reader
}

func NewFastaScanner(reader io.Reader) FastaScanner {
	return FastaScanner{reader: bufio.NewReader(reader)}
}

// Scan will return the next token, byte literal of the
// non-interleaved data, and the length of the token's
// value
func (f FastaScanner) Scan() (Token, []byte, int) {
	ch, size, err := f.reader.ReadRune()
	if err != nil {
		return EOF, []byte{}, 0
	}
	switch {
	case ch == '>':
		lit, length := f.scanSequenceId()
		return SEQUENCE_ID, lit, length
	case size > 1:
		return INVALID, []byte{}, 0
	case scanner.IsSequenceData(ch):
		f.reader.UnreadRune()
		lit, length, err := scanner.ScanSequenceData(f.reader)
		if err != nil && err != io.EOF {
			return INVALID, []byte{}, 0
		}
		return SEQUENCE_DATA, lit, length
	case scanner.IsWhitespace(ch):
		f.reader.UnreadByte()
		lit, length, err := scanner.ScanWhitespace(f.reader)
		if err != io.EOF {
			return INVALID, []byte{}, 0
		}
		return WHITESPACE, lit, length
	default:
		return EOF, []byte{}, 0
	}
}

// scanSequenceId from the fasta file.  This is the part following
// a '>'.  I cheat by using the read string til \n,
// it returns the length of the id string
func (f FastaScanner) scanSequenceId() ([]byte, int) {
	line, err := f.reader.ReadSlice('\n')
	if err != nil {
		return line, 0
	}
	return line[:len(line)-1], len(line) - 1
}
