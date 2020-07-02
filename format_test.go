package syslogp

import (
	"fmt"
	"io/ioutil"
	"log/syslog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatRFC5424WithFields(t *testing.T) {
	expectedFieldsOutput := "[context@123 key1=\"value1\" key2=\"value2\" key3=\"value3\"]"
	testRFC5424WithFields(t, expectedFieldsOutput, Fields{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})
}

func TestFormatRFC5424WithoutFields(t *testing.T) {
	expectedFieldsOutput := "-"
	testRFC5424WithFields(t, expectedFieldsOutput, nil)
}

func testRFC5424WithFields(t *testing.T, expectedFieldsString string, fields Fields) {
	expectedOutput := fmt.Sprintf("<158>1 1985-04-12T23:20:50Z foo.bar testing_app 123 errors_trace %v Im a log! :-)\n", expectedFieldsString)

	now, err := time.Parse(time.RFC3339, "1985-04-12T23:20:50.52Z")
	assert.NoError(t, err)

	formatter, err := NewRFC5424Formatter(
		syslog.LOG_LOCAL3,
		syslog.LOG_INFO,
		now,
		"foo.bar",
		"testing_app",
		123,
		"errors_trace",
		fields,
		"Im a log! :-)",
	)
	assert.NoError(t, err)

	result, err := ioutil.ReadAll(formatter)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(result))
}
func TestFormatRFC3164(t *testing.T) {
	expectedOutput := "<158>Jan 29 18:39:56 foo.bar testing_app[123]: Im a log! :-)\n"

	now, err := time.Parse(time.Stamp, "Jan 29 18:39:56")
	assert.NoError(t, err)

	formatter, err := NewRFC3164Formatter(
		syslog.LOG_LOCAL3,
		syslog.LOG_INFO,
		now,
		"foo.bar",
		"testing_app",
		123,
		"Im a log! :-)",
	)
	assert.NoError(t, err)

	result, err := ioutil.ReadAll(formatter)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(result))
}
