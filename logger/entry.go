package logger

import "time"

//Entry is the structure used for logging.
type Entry struct {
	message string
	fields  []Field
	level   Level
	time    time.Time
}

// NewFakeEntry creates a fake entry for testing purposes.
func NewFakeEntry() *Entry {
	f := NewField("user_id", 12345)
	return &Entry{
		"foo message",
		[]Field{f},
		LevelDebug,
		time.Now(),
	}
}
