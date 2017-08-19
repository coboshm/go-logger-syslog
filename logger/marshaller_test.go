package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogstashMarshaller(t *testing.T) {
	entry := NewFakeEntry()
	expectedOutput := []byte(
		fmt.Sprintf("{\"@timestamp\":\"%s\",\"@version\":1,\"app_server_name\":\"%s\",\"application\":\"api\",\"channel\":\"my-chan\",\"env\":\"test\",\"message\":\"foo message\",\"severity\":\"debug\",\"user_id\":\"12345\"}\n",
			entry.time.Format(time.RFC3339),
			hostname,
		),
	)
	marshaller := newLogstashMarshaller("my-chan", "api", "test")

	marshaledData, err := marshaller.Marshal(entry)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, marshaledData)
}

func TestSyslogToLogstashMarshaller(t *testing.T) {
	entry := NewFakeEntry()
	expectedOutput := []byte(
		fmt.Sprintf(
			"@cee: {\"@timestamp\":\"%s\",\"@version\":1,\"app_server_name\":\"%s\",\"application\":\"api\",\"channel\":\"my-chan\",\"env\":\"test\",\"message\":\"foo message\",\"severity\":\"debug\",\"user_id\":\"12345\"}\n",
			entry.time.Format(time.RFC3339),
			hostname,
		),
	)

	marshaller := newSyslogToLogstashMarshaller("my-chan", "api", "test")

	marshaledData, err := marshaller.Marshal(entry)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, marshaledData)
}
