package logsplitter

import (
	"fmt"
	"regexp"
	"strings"
)

// StringSplitter the the interface that exposes string splitting functionality
type StringSplitter interface {
	Split(input string) ([]string, error)
}

// DelimeterStringSplitter splits a string based on a delimeter
type DelimeterStringSplitter struct {
	delim string
}

// NewDelimeterStringSplitter creates a new DelimeterStringSplitter
func NewDelimeterStringSplitter(delim string) (*DelimeterStringSplitter, error) {
	var splitter *DelimeterStringSplitter
	if delim == "" {
		return splitter, fmt.Errorf("%q is not a valid delimeter", delim)
	}

	splitter = &DelimeterStringSplitter{
		delim: delim,
	}

	return splitter, nil
}

// Split splits a string bases on a delimeter
func (s *DelimeterStringSplitter) Split(input string) ([]string, error) {
	return strings.Split(input, s.delim), nil
}

// RegexStringSplitter splits a string with a regular expression that contains submatches for each part of a log line
type RegexStringSplitter struct {
	regex   *regexp.Regexp
	pattern string
}

// NewRegexStringSplitter creates a new RegexStringSplitter
func NewRegexStringSplitter(pattern string) (*RegexStringSplitter, error) {
	var splitter RegexStringSplitter
	r, err := regexp.Compile(pattern)
	if err != nil {
		return &splitter, fmt.Errorf("unable to parse regex - %v", err)
	}

	splitter = RegexStringSplitter{
		regex:   r,
		pattern: pattern,
	}

	return &splitter, nil
}

// Split splits a string into parts with a regular expression that contains submatches
func (s *RegexStringSplitter) Split(input string) ([]string, error) {
	result := s.regex.FindStringSubmatch(input)

	if len(result) > 0 {
		return result[1:], nil
	} else {
		return result, nil
	}
}
