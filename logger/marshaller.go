package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

// logInfo  ...
type logInfo struct {
	channel     string
	application string
	environment string
	hostname    string
}

// Marshaller marshals
type Marshaller interface {
	// Marshal returns the info encoded to be readable by humans
	Marshal(entry *Entry) ([]byte, error)
}

type noopMarshaller struct{}

func newNOOPMarshaller() Marshaller {
	return &noopMarshaller{}

}

// Marshal  ...
func (n *noopMarshaller) Marshal(entry *Entry) ([]byte, error) {
	return nil, nil
}

type humanMarshaller struct {
	logInfo
}

// newHumanMarshaller ...
func newHumanMarshaller(channel, application, environment string) Marshaller {
	return &humanMarshaller{
		logInfo{
			channel:     channel,
			application: application,
			environment: environment,
			hostname:    hostname,
		},
	}
}

// Marshal returns the info encoded to be readable by humans.
func (m *humanMarshaller) Marshal(entry *Entry) ([]byte, error) {
	separator := ", "
	var buffer bytes.Buffer

	lvl, _ := entry.level.String()
	buffer.WriteString(fmt.Sprintf("[%v]%v", entry.time.Format(time.RFC3339), separator))
	buffer.WriteString(fmt.Sprintf("%v.%v%v", m.application, m.environment, separator))
	buffer.WriteString(fmt.Sprintf("%v.%v%v", m.channel, lvl, separator))
	buffer.WriteString(fmt.Sprintf("%v%v", entry.message, separator))
	buffer.WriteString("[")
	// rest
	for i, field := range entry.fields {
		buffer.WriteString(fmt.Sprintf("%v:%v", field.Key, field.Value))
		if i != len(entry.fields)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]\n")

	return buffer.Bytes(), nil
}

// jSONMarshaller implements the Marshaller interface using the json marshaller from the std lib
type jSONMarshaller struct{}

func newJSONMarshaller() Marshaller {
	return &jSONMarshaller{}
}

// Marshal returns the info encoded to be readable by humans
func (m *jSONMarshaller) Marshal(entry *Entry) ([]byte, error) {
	return json.Marshal(entry)
}

// syslogToLogstashMarshaller marshalls the data to a logstash-compatible JSON adding a cee suffix.
// See http://www.rsyslog.com/doc/mmjsonparse.html
type syslogToLogstashMarshaller struct {
	logInfo
}

// Marshal returns the info encoded to be readable by humans
func (m *syslogToLogstashMarshaller) Marshal(entry *Entry) ([]byte, error) {
	logstashMarshaller := newLogstashMarshaller(m.logInfo.channel, m.logInfo.application, m.logInfo.environment)
	jsonEncoded, err := logstashMarshaller.Marshal(entry)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.Write([]byte("@cee: "))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(jsonEncoded)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// newSyslogToLogstashMarshaller is the constructor of the concrete type.
func newSyslogToLogstashMarshaller(channel, application, environment string) Marshaller {
	return &syslogToLogstashMarshaller{
		logInfo{
			channel:     channel,
			application: application,
			environment: environment,
			hostname:    hostname,
		},
	}
}

// logstashMarshaller marshalls the data to a logstash-compatible JSON.
type logstashMarshaller struct {
	logInfo
}

// newLogstashMarshaller is the constructor of the concrete type.
func newLogstashMarshaller(channel, application, environment string) Marshaller {
	return &logstashMarshaller{
		logInfo{
			channel:     channel,
			application: application,
			environment: environment,
			hostname:    hostname,
		},
	}
}

// Marshal returns the info encoded in the logstash format (JSON with special fields)
func (l *logstashMarshaller) Marshal(entry *Entry) ([]byte, error) {

	// these are words that are used by logstash so we have to change its name when
	// they are processed.
	var reservedWords = map[string]struct{}{
		"error": struct{}{},
		"type":  struct{}{},
	}

	data := make(map[string]interface{})
	// logstash ones
	data["@version"] = 1
	data["@timestamp"] = entry.time.Format(time.RFC3339)
	severity, _ := entry.level.String()

	data["severity"] = severity
	data["message"] = entry.message
	data["app_server_name"] = l.hostname
	data["channel"] = l.channel
	data["application"] = l.application
	data["env"] = l.environment
	// rest
	for _, field := range entry.fields {
		value := fmt.Sprintf("%v", field.Value)
		if _, ok := reservedWords[field.Key]; ok {
			data[fmt.Sprintf("%sx", field.Key)] = value
		} else {
			data[field.Key] = value
		}
	}
	encodedData, err := json.Marshal(data)
	// as the std logger does
	if len(encodedData) == 0 || encodedData[len(encodedData)-1] != '\n' {
		encodedData = append(encodedData, '\n')
	}
	return encodedData, err
}
