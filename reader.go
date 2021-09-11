package logsplitter

import (
	"bufio"
	"io"
)

// ParseReader reads strings, parses it into seperate parts and return the seperate parts as a logsplitter.Fields
type ParseReader struct {
	scanner *bufio.Scanner
	parser  *StringParser
}

// NewParseReader creates a new ParseReader
func NewParseReader(file io.Reader, parser *StringParser) *ParseReader {
	reader := ParseReader{
		scanner: bufio.NewScanner(file),
		parser:  parser,
	}
	return &reader
}

// Read reads one log file line, parses it into seperate parts and return the seperate parts as a logsplitter.Fields
func (r *ParseReader) Read() (Fields, error) {
	var result Fields
	ok := r.scanner.Scan()
	if !ok {
		return result, io.EOF
	}

	result, err := r.parser.Parse(r.scanner.Text())
	if err != nil {
		return result, err
	}
	return result, nil
}
