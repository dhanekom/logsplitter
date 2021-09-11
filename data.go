package logsplitter

// Field represents a part (section) of a log line
type Field struct {
	Name  string
	Value string
}

// Fields represents the data of one log line that has been split into parts
type Fields []*Field
