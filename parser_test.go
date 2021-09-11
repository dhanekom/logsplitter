package logsplitter

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type FailStringSplitter struct {
}

func (s FailStringSplitter) Split(input string) ([]string, error) {
	var result []string
	return result, errors.New("some error")
}

func TestDelimParse(t *testing.T) {
	tests := []struct {
		desc                string
		input               string
		splitDelim          string
		concatDelim         string
		columnOnly          bool
		columnIndex         int
		output              string
		result              bool
		failStringSplitting bool
	}{
		{"No delims", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "", "", false, 0, "", false, false},
		{"Row same delims", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "|", "|", false, 0, "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", true, false},
		{"Row different delims", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "|", ",", false, 0, "2021/08/30 19:41:15.740,INFO,2,1553,Starting SKU Refresh", true, false},
		{"Row fail string splitting", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "|", ",", false, 0, "fresh", false, true},
		{"Column valid value", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "|", ",", true, 3, "1553", true, false},
		{"Column invalid index", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", "|", ",", true, 5, "", false, false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			var sb strings.Builder
			var stringSplitter StringSplitter
			var err error
			if test.failStringSplitting {
				stringSplitter = FailStringSplitter{}
			} else {
				stringSplitter, err = NewDelimeterStringSplitter(test.splitDelim)
			}
			if err != nil {
				if test.result {
					t.Error(err)
				}
			} else {
				parser := NewParser(stringSplitter)
				reader := NewParseReader(strings.NewReader(test.input), parser)

				for {
					fields, err := reader.Read()
					if err != nil {
						if err == io.EOF {
							break
						}

						if test.result {
							t.Errorf("expected to parse string successully but got error - %v", err)
							return
						}
					} else {
						if test.columnOnly {
							if len(fields)-1 < test.columnIndex && test.result {
								t.Errorf("attempting to access row field index %d while max index is %d", test.columnIndex, len(fields)-1)
							} else {
								if test.result && fields[test.columnIndex].Value != test.output {
									t.Errorf("expected %s, got %s", test.output, fields[test.columnIndex].Value)
								}
							}
						} else {
							for _, field := range fields {
								if sb.String() != "" {
									sb.WriteString(test.concatDelim)
								}
								sb.WriteString(field.Value)
							}

							if test.output != sb.String() {
								t.Errorf("expected %s, got %s", test.output, sb.String())
							}
						}
					}
				}
			}
		})
	}
}

func TestRegexParse(t *testing.T) {

	const pattern = `^([0-9\-/ :\.]{0,23})\|(\w+)\|(\d+)\|(\-?\d+)\|(.+)$`
	tests := []struct {
		desc        string
		input       string
		pattern     string
		concatDelim string
		columnOnly  bool
		columnIndex int
		output      string
		result      bool
	}{
		{"Invalid regex", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", `(`, "|", false, 0, "", false},
		{"Blank regex", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", ``, "|", false, 0, "", true},
		{"No matches", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", `(123)`, "|", false, 0, "", true},
		{"Row same delims", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", pattern, "|", false, 0, "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", true},
		{"Row different delims", "2021/08/30 19:41:15.740|INFO|2|-99|Starting SKU Refresh", pattern, ",", false, 0, "2021/08/30 19:41:15.740,INFO,2,-99,Starting SKU Refresh", true},
		{"Column valid value", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", pattern, ",", true, 3, "1553", true},
		{"Column invalid index", "2021/08/30 19:41:15.740|INFO|2|1553|Starting SKU Refresh", pattern, ",", true, 5, "", false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			var sb strings.Builder
			stringSplitter, err := NewRegexStringSplitter(test.pattern)
			if err != nil {
				if test.result {
					t.Error(err)
				}
			} else {
				parser := NewParser(stringSplitter)
				reader := NewParseReader(strings.NewReader(test.input), parser)

				for {
					fields, err := reader.Read()
					if err != nil {
						if err == io.EOF {
							break
						}

						if test.result {
							t.Errorf("expected to parse string successully but got error - %v", err)
							return
						}
					} else {
						if test.columnOnly {
							if len(fields)-1 < test.columnIndex && test.result {
								t.Errorf("attempting to access row field index %d while max index is %d", test.columnIndex, len(fields)-1)
							} else {
								if test.result && fields[test.columnIndex].Value != test.output {
									t.Errorf("expected %s, got %s", test.output, fields[test.columnIndex].Value)
								}
							}
						} else {
							for _, field := range fields {
								if sb.String() != "" {
									sb.WriteString(test.concatDelim)
								}
								sb.WriteString(field.Value)
							}

							if test.output != sb.String() {
								t.Errorf("expected %s, got %s", test.output, sb.String())
							}
						}
					}
				}
			}
		})
	}
}
