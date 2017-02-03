package syslogp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/syslog"
	"strings"
	"time"
)

// Right now there is only version 1.
const syslogVersion = 1

// Fields for the RFC5424.
type Fields map[string]interface{}

type rfc5424Formatter struct {
	priority syslog.Priority
	time     time.Time
	hostname string
	app      string
	pid      int
	msgID    string
	fields   Fields
	msg      string
	buf      bytes.Buffer
}

// NewRFC5424Formatter returns a formatter for this specific RFC.
// See https://tools.ietf.org/html/rfc5424
func NewRFC5424Formatter(facility syslog.Priority, severity syslog.Priority, time time.Time, hostname string, app string, pid int, msgID string, fields Fields, msg string) (io.Reader, error) {
	priority := facility | severity
	if priority < 0 || priority > syslog.LOG_LOCAL7|syslog.LOG_DEBUG {
		return nil, errors.New("invalid syslog priority")
	}

	var buf bytes.Buffer

	return &rfc5424Formatter{
		priority, time, hostname,
		app, pid, msgID,
		fields, msg,
		buf,
	}, nil
}

func (f *rfc5424Formatter) Read(p []byte) (n int, err error) {
	if f.buf.Len() > 0 {
		return 0, io.EOF
	}

	fields := "-"
	var i int
	if len(f.fields) > 0 {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("[context@%v ", f.pid))
		for key, value := range f.fields {
			buf.WriteString(fmt.Sprintf(`%v="%v"`, key, value))
			if i != len(f.fields)-1 {
				buf.WriteString(" ")
			}
			i++
		}
		buf.WriteString("]")

		fields = buf.String()
	}

	// ensure it ends in a \n
	nl := ""
	if !strings.HasSuffix(f.msg, "\n") {
		nl = "\n"
	}

	// <priority>VERSION ISOTIMESTAMP HOSTNAME APPLICATION PID MESSAGEID STRUCTURED-DATA MSG
	_, err = fmt.Fprintf(
		&f.buf,
		"<%d>%v %v %v %v %v %v %v %v%v",
		f.priority,
		syslogVersion,
		f.time.Format(time.RFC3339),
		f.hostname,
		f.app,
		f.pid,
		f.msgID,
		fields,
		f.msg,
		nl,
	)

	if err != nil {
		return 0, err
	}

	copy(p, f.buf.Bytes())

	return f.buf.Len(), nil
}

type rfc3164Formatter struct {
	priority syslog.Priority
	time     time.Time
	hostname string
	app      string
	pid      int
	msg      string
	buf      bytes.Buffer
}

// NewRFC3164Formatter returns a formatter for this specific RFC.
// See https://tools.ietf.org/html/rfc3164
func NewRFC3164Formatter(facility syslog.Priority, severity syslog.Priority, time time.Time, hostname string, app string, pid int, msg string) (io.Reader, error) {
	priority := facility | severity
	if priority < 0 || priority > syslog.LOG_LOCAL7|syslog.LOG_DEBUG {
		return nil, errors.New("invalid syslog priority")
	}

	var buf bytes.Buffer

	return &rfc3164Formatter{
		priority, time, hostname,
		app, pid, msg,
		buf,
	}, nil
}

func (f *rfc3164Formatter) Read(p []byte) (n int, err error) {
	if f.buf.Len() > 0 {
		return 0, io.EOF
	}

	// ensure it ends in a \n
	nl := ""
	if !strings.HasSuffix(f.msg, "\n") {
		nl = "\n"
	}

	//<PRI>TIMESTAMP HOSTNAME TAG[PID]: MSG
	_, err = fmt.Fprintf(
		&f.buf,
		"<%d>%s %s %s[%d]: %s%s",
		f.priority,
		f.time.Format(time.Stamp),
		f.hostname,
		f.app,
		f.pid,
		f.msg,
		nl,
	)

	if err != nil {
		return 0, err
	}

	copy(p, f.buf.Bytes())

	return f.buf.Len(), nil
}
