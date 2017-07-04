// +build windows plan9

package logger

import "fmt"

// DefaultConfig is
func DefaultConfig() *Config {
	result := &Config{
		LogLevel:       "INFO",
		EnableSyslog:   false,
		SyslogFacility: "",
	}
	return result
}

// NewLogger is used to construct a new Syslogger
func NewLogger(p Priority, facility, tag string) (Syslogger, error) {
	return nil, fmt.Errorf("Platform does not support syslog")
}

// DialLogger is used to construct a new Syslogger that establishes connection to remote syslog server
func DialLogger(network, raddr string, p Priority, facility, tag string) (Syslogger, error) {
	return nil, fmt.Errorf("Platform does not support syslog")
}
