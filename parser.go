package logsplitter

import "fmt"

// StringParser parses a log file by splitting each log file line into seperate parts and returns the parts as a logsplitter.Fields
type StringParser struct {
	splitter StringSplitter
}

// NewParser creates a new StringParser
func NewParser(splitter StringSplitter) *StringParser {
	return &StringParser{
		splitter: splitter,
	}
}

// Parse parses a log file line into seperate parts and returns the parts as a logsplitter.Fields
func (p *StringParser) Parse(input string) (Fields, error) {
	var result Fields

	values, err := p.splitter.Split(input)
	if err != nil {
		return result, fmt.Errorf("unable to split strings while parsing - %v", err)
	}

	for _, value := range values {
		field := Field{
			Name:  "",
			Value: value,
		}
		result = append(result, &field)
	}

	return result, nil
}
