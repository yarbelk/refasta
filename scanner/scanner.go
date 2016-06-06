package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"
)

type InvalidChar error

var eof = rune(0)

// ScanWhitespace returns a contigious block of whitespaces, and return
// literal bytes for the token, length of the string and error
func ScanWhitespace(reader *bufio.Reader) (lit []byte, length int, err error) {
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.Bytes(), readErr
			break scanLoop
		case size > 1:
			err = InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case IsWhitespace(ch):
			length++
			buf.WriteRune(ch)
			continue scanLoop
		default:
			reader.UnreadRune()
			lit, err = buf.Bytes(), nil
			break scanLoop
		}
	}
	return
}

// IsWhitespace returns true if character is whitespace
func IsWhitespace(c rune) bool {
	return unicode.IsSpace(c)
}

// IsSequenceData returns true if the character is a letter/number
// Number is because Morphology data generally is 0..9, while
// DNA/RNA/Proteins are letters.  All are valid
func IsSequenceData(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '-' || c == '?'
}

// ScanSequenceData will return a string of sequence data, removing all
// new line characters
func ScanSequenceData(reader *bufio.Reader) (lit []byte, length int, alphabet map[rune]bool, err error) {
	alphabet = make(map[rune]bool)
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.Bytes(), readErr
			break scanLoop
		case size > 1:
			err = InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case IsSequenceData(ch):
			length++
			buf.WriteRune(ch)
			alphabet[ch] = true
			continue scanLoop
		case ch == '[':
			subSeq, subLen, alpha, err := scanSequenceDataGroup(reader)
			for k, _ := range alpha {
				alphabet[k] = true
			}
			if err != nil {
				return subSeq, subLen, alphabet, err
			}
			buf.Write(subSeq)
			length = length + subLen
			continue scanLoop
		case IsWhitespace(ch):
			// skip whitespace
			continue scanLoop
		case ch == '>':
			reader.UnreadRune()
			lit, err = buf.Bytes(), nil
			break scanLoop
		case ch == eof:
			lit, err = buf.Bytes(), nil
			break scanLoop
		default:
			lit, err = buf.Bytes(), InvalidChar(fmt.Errorf("No idea what this is '%s' in the char stream", string(ch)))
			break scanLoop
		}
	}
	return
}

// scanSequenceDataGroup handles the scanning of [ATGA] like sequence
// structures.
func scanSequenceDataGroup(reader *bufio.Reader) (lit []byte, length int, alphabet map[rune]bool, err error) {
	alphabet = make(map[rune]bool)
	buf := bytes.Buffer{}

scanLoop:
	for {
		ch, size, readErr := reader.ReadRune()
		switch {
		case err != nil:
			lit, err = buf.Bytes(), readErr
			break scanLoop
		case size > 1:
			err = InvalidChar(fmt.Errorf("Invalid Char in stream, %s", ch))
			break scanLoop
		case IsWhitespace(ch):
			continue scanLoop
		case IsSequenceData(ch):
			buf.WriteRune(ch)
			alphabet[ch] = true
			continue scanLoop
		case ch == ']':
			length++
			buf.WriteRune(ch)
			lit, err = buf.Bytes(), nil
			break scanLoop
		case ch == '>', ch == eof:
			return []byte{}, 0, nil, InvalidChar(fmt.Errorf("Unbalanced [] in sequence data: postion %d", length))
		default:
			lit, err = buf.Bytes(), InvalidChar(fmt.Errorf("No idea what this is '%s' in the char stream", ch))
			break scanLoop
		}
	}
	return
}
