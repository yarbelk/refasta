package sequence

import "fmt"

type ErrNo int

const (
	UNKNOWN ErrNo = iota
	MISSMATCHED_SEQUENCE_LENGTHS
	BAD_FORMAT
)

// InvalidSequence is an error type that (will) hold useful data about
// bad data when preparing to serialize to a format
type InvalidSequence struct {
	Message string
	Details string
	Errno   ErrNo
}

type FormatError struct {
	Message string
	Details string
	Errno   ErrNo
}

// Error for the error interface
func (e InvalidSequence) Error() string {
	return fmt.Sprintf("InvalidSequence: %s\nDetails: %s", e.Message, e.Details)
}

func (e FormatError) Error() string {
	return fmt.Sprintf("FormatError: %s\nDetails: %s", e.Message, e.Details)
}
