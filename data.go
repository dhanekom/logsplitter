package logsplitter

// Field represents a part (section) of a log line
type Field struct {
	name  string
	value string
}

// Fields represents the data of one log line that has been split into parts
type Fields []*Field

// Name returns the string value of a field
func (f Field) Value() string {
	return f.value
}

// Name returns the name of a field
func (f Field) Name() string {
	return f.name
}
