package formats

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

const (
	EOF Token = iota
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

func (f FastaScanner) Scan() (Token, string) {
	ch, size, err := f.reader.ReadRune()
	if err != nil {
		return EOF, ""
	}
	switch {
	case ch == '>':
		return SEQUENCE_ID, f.scanSequenceId()
	case size > 1:
		return INVALID, ""
	case isSequenceData(ch):
		f.reader.UnreadRune()
		lit, err := f.scanSequenceData()
		if err != nil && err != io.EOF {
			return INVALID, ""
		}
		return SEQUENCE_DATA, lit
	case isWhitespace(ch):
		f.reader.UnreadByte()
		lit, err := f.scanWhitespace()
		if err != io.EOF {
			return INVALID, ""
		}
		return WHITESPACE, lit
	default:
		return EOF, ""
	}
}

// scanSequenceId from the fasta file.  This is the part following
// a '>'.  I cheat by using the read string til \n
func (f FastaScanner) scanSequenceId() string {
	line, err := f.reader.ReadString('\n')
	if err != nil {
		return line
	}
	return line[:len(line)-1]
}

// scanWhitespace returns a contigious block of whitespaces
func (f FastaScanner) scanWhitespace() (lit string, err error) {
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := f.reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.String(), readErr
			break scanLoop
		case size > 1:
			lit, err = "", InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case isWhitespace(ch):
			buf.WriteRune(ch)
			continue scanLoop
		default:
			f.reader.UnreadRune()
			lit, err = buf.String(), nil
			break scanLoop
		}
	}
	return
}

// scanSequenceData will return a string of sequence data, removing all
// new line characters
func (f FastaScanner) scanSequenceData() (lit string, err error) {
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := f.reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.String(), readErr
			break scanLoop
		case size > 1:
			lit, err = "", InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case isWhitespace(ch):
			continue scanLoop
		case ch == '>':
			f.reader.UnreadRune()
			lit, err = buf.String(), nil
			break scanLoop
		case isSequenceData(ch):
			buf.WriteRune(ch)
			continue scanLoop
		case ch == eof:
			lit, err = buf.String(), nil
			break scanLoop
		default:
			lit, err = buf.String(), InvalidChar(fmt.Errorf("No idea what this is '%s' in the char stream", ch))
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
