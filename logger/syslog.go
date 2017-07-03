package logger

import (
	"bytes"
)

// Priority maps to the syslog priority levels
type Priority int

// consts
const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

var levelPriority = map[string]Priority{
	"TRACE": LOG_DEBUG,
	"DEBUG": LOG_INFO,
	"INFO":  LOG_NOTICE,
	"WARN":  LOG_WARNING,
	"ERR":   LOG_ERR,
	"CRIT":  LOG_CRIT,
}

// SyslogWrapper is used to cleanup log messages before
// writing them to a Syslogger
type SyslogWrapper struct {
	l      Syslogger
	filter *LevelFilter
}

// Syslogger interface is used to write log messages to syslog
type Syslogger interface {
	WriteLevel(Priority, []byte) error
	Write([]byte) (int, error)
	Close() error
}

// Write is used to implement io.Writer
func (s *SyslogWrapper) Write(p []byte) (int, error) {
	if !s.filter.Check(p) {
		return 0, nil
	}

	var level string
	afterLevel := p
	x := bytes.IndexByte(p, '[')
	if x >= 0 {
		y := bytes.IndexByte(p[x:], ']')
		if y >= 0 {
			level = string(p[x+1 : x+y])
			afterLevel = p[x+y+2:]
		}
	}

	priority, ok := levelPriority[level]
	if !ok {
		priority = LOG_NOTICE
	}

	err := s.l.WriteLevel(priority, afterLevel)
	return len(p), err
}
