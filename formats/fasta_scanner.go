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
	case isSequenceData(ch):
		f.reader.UnreadRune()
		lit, length, err := f.scanSequenceData()
		if err != nil && err != io.EOF {
			return INVALID, []byte{}, 0
		}
		return SEQUENCE_DATA, lit, length
	case isWhitespace(ch):
		f.reader.UnreadByte()
		lit, length, err := f.scanWhitespace()
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

// scanWhitespace returns a contigious block of whitespaces, and return
// literal bytes for the token, length of the string and error
func (f FastaScanner) scanWhitespace() (lit []byte, length int, err error) {
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
			length++
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
func (f FastaScanner) scanSequenceData() (lit []byte, length int, err error) {
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
		case ch == '[':
			subSeq, subLen, err := f.scanSequenceDataGroup()
			if err != nil {
				return subSeq, subLen, err
			}
			buf.Write(subSeq)
			length = length + subLen
			continue scanLoop
		case isWhitespace(ch):
			continue scanLoop
		case ch == '>':
			f.reader.UnreadRune()
			lit, err = buf.Bytes(), nil
			break scanLoop
		case isSequenceData(ch):
			length++
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

// scanSequenceDataGroup handles the scanning of [ATGA] like sequence
// structures.
func (f FastaScanner) scanSequenceDataGroup() (lit []byte, length int, err error) {
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
		case isSequenceData(ch):
			buf.WriteRune(ch)
			continue scanLoop
		case ch == ']':
			length++
			buf.WriteRune(ch)
			lit, err = buf.Bytes(), nil
			break scanLoop
		case ch == '>', ch == eof:
			return []byte{}, 0, InvalidChar(fmt.Errorf("Unbalanced [] in sequence data: postion %d", length))
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
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '-'
}
