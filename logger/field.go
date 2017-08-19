package logger

// Field is a key/value pair associated to a log.
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new field.
func NewField(key string, value interface{}) Field {
	return Field{
		key,
		value,
	}
}
