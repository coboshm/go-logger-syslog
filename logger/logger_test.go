package logger

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebugLoggerMessage(t *testing.T) {
	marshaller := marshallerMock()
	logField := NewField("user_id", 12345)
	logMessage := "this is a debug message"

	writer := writerMock(t, logMessage, marshaller)
	logger := newLog(marshaller, writer, LevelDebug)

	logger.Debug(logMessage, logField)
	assert.Equal(t, 1, writer.WriteTotalCalls())
}

func TestDoesNotLogForUnexpectedLevel(t *testing.T) {
	marshaller := marshallerMock()
	logField := NewField("user_id", 12345)
	logMessage := "foo message"

	writer := writerMock(t, logMessage, marshaller)
	logger := newLog(marshaller, writer, LevelInfo)

	logger.Debug(logMessage, logField)
	assert.Equal(t, 0, writer.WriteTotalCalls())
}

func TestCreatesAStdOutLoggerFromDSN(t *testing.T) {
	_, err := NewLoggerFromDSN("stdout://?level=debug", "app", "test")
	assert.NoError(t, err)
}

func TestCreatesADiscardAllLoggerFromDSN(t *testing.T) {
	_, err := NewLoggerFromDSN("discardall://?level=debug", "app", "test")
	assert.NoError(t, err)
}

func TestErrorWhenDSNIsInvalid(t *testing.T) {
	_, err := NewLoggerFromDSN("foo://?level=debug", "app", "test")
	assert.Error(t, err)
}

func TestErrorWhenLevelIsInvalid(t *testing.T) {
	_, err := NewLoggerFromDSN("discardall://?level=x", "app", "test")
	assert.Error(t, err)
}

func TestCreatesALoggerWithDefaultLevelWhenItsMissing(t *testing.T) {
	_, err := NewLoggerFromDSN("discardall://", "app", "test")
	assert.NoError(t, err)
}

func marshallerMock() Marshaller {
	marshaller := NewMarshallerMock()
	marshaller.MarshalFunc = func(e *Entry) ([]byte, error) {
		return []byte(e.message), nil
	}

	return marshaller
}

func writerMock(t *testing.T, expectedMessage string, marshaller Marshaller) *WriterMock {
	mockWriter := NewWriterMock()
	mockWriter.WriteFunc = func(data []byte) (int, error) {
		entry := &Entry{
			message: expectedMessage,
			fields:  nil,
			level:   0,
			time:    time.Now(),
		}
		expectedData, _ := marshaller.Marshal(entry)
		assert.Equal(t, expectedData, data)

		return 0, nil
	}

	return mockWriter
}
