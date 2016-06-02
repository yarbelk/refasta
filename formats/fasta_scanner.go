package formats

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

const (
	UNSTARTED Token = -1
	EOF       Token = iota
	SEQUENCE_ID
	SEQUENCE_DATA
	WHITESPACE
	INVALID
)

type InvalidChar error

var eof = rune(0)

type Token int

type FastaScanner struct {
	reader *bufio.Reader
}

func NewFastaScanner(reader io.Reader) FastaScanner {
	return FastaScanner{reader: bufio.NewReader(reader)}
}

func (f FastaScanner) Scan() (Token, []byte) {
	ch, size, err := f.reader.ReadRune()
	if err != nil {
		return EOF, []byte{}
	}
	switch {
	case ch == '>':
		return SEQUENCE_ID, f.scanSequenceId()
	case size > 1:
		return INVALID, []byte{}
	case isSequenceData(ch):
		f.reader.UnreadRune()
		lit, err := f.scanSequenceData()
		if err != nil && err != io.EOF {
			return INVALID, []byte{}
		}
		return SEQUENCE_DATA, lit
	case isWhitespace(ch):
		f.reader.UnreadByte()
		lit, err := f.scanWhitespace()
		if err != io.EOF {
			return INVALID, []byte{}
		}
		return WHITESPACE, lit
	default:
		return EOF, []byte{}
	}
}

// scanSequenceId from the fasta file.  This is the part following
// a '>'.  I cheat by using the read string til \n
func (f FastaScanner) scanSequenceId() []byte {
	line, err := f.reader.ReadSlice('\n')
	if err != nil {
		return line
	}
	return line[:len(line)-1]
}

// scanWhitespace returns a contigious block of whitespaces
func (f FastaScanner) scanWhitespace() (lit []byte, err error) {
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := f.reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.Bytes(), readErr
			break scanLoop
		case size > 1:
			err = InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case isWhitespace(ch):
			buf.WriteRune(ch)
			continue scanLoop
		default:
			f.reader.UnreadRune()
			lit, err = buf.Bytes(), nil
			break scanLoop
		}
	}
	return
}

// scanSequenceData will return a string of sequence data, removing all
// new line characters
func (f FastaScanner) scanSequenceData() (lit []byte, err error) {
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := f.reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.Bytes(), readErr
			break scanLoop
		case size > 1:
			err = InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case isWhitespace(ch):
			continue scanLoop
		case ch == '>':
			f.reader.UnreadRune()
			lit, err = buf.Bytes(), nil
			break scanLoop
		case isSequenceData(ch):
			buf.WriteRune(ch)
			continue scanLoop
		case ch == eof:
			lit, err = buf.Bytes(), nil
			break scanLoop
		default:
			lit, err = buf.Bytes(), InvalidChar(fmt.Errorf("No idea what this is '%s' in the char stream", ch))
			break scanLoop
		}
	}
	return
}

// isWhitespace returns true if character is whitespace
func isWhitespace(c rune) bool {
	return unicode.IsSpace(c)
}

// isSequenceData returns true if the character is a letter/number
// Number is because Morphology data generally is 0..9, while
// DNA/RNA/Proteins are letters.  All are valid
func isSequenceData(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c)
}
