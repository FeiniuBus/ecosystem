package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Config is used to set up logging.
type Config struct {
	LogLevel       string
	EnableSyslog   bool
	SyslogFacility string
}

// Setup is used to perform setup of serveral logging objects
func setup(config *Config) (*LevelFilter, *GatedWriter, *LogWriter, io.Writer, bool) {
	logGate := &GatedWriter{
		Writer: os.Stdout,
	}

	logFilter := NewLevelFilter()
	logFilter.MinLevel = LogLevel(strings.ToUpper(config.LogLevel))
	logFilter.Writer = logGate
	if !ValidateLevelFilter(logFilter.MinLevel, logFilter) {
		fmt.Fprintf(os.Stderr, "Invalid log level: %s. Valid log levels are: %v", logFilter.MinLevel, logFilter.Levels)
		return nil, nil, nil, nil, false
	}

	var syslog io.Writer
	if config.EnableSyslog {
		retries := 12
		delay := 5 * time.Second
		for i := 0; i <= retries; i++ {
			l, err := NewLogger(LOG_NOTICE, config.SyslogFacility, "FEINIUBUS")
			if err == nil {
				syslog = &SyslogWrapper{l, logFilter}
				break
			}

			fmt.Fprintf(os.Stderr, "Syslog setup error: %v", err)
			if i == retries {
				timeout := time.Duration(retries) * delay
				fmt.Fprintf(os.Stderr, "Syslog setup did not succeed within timeout (%s)", timeout.String())
				return nil, nil, nil, nil, false
			}

			fmt.Fprintf(os.Stderr, "Retrying syslog setup in %s...", delay.String())
			time.Sleep(delay)
		}
	}

	logWriter := NewLogWriter(512)
	var logOutput io.Writer
	if syslog != nil {
		logOutput = io.MultiWriter(logFilter, logWriter, syslog)
	} else {
		logOutput = io.MultiWriter(logFilter, logWriter)
	}
	return logFilter, logGate, logWriter, logOutput, true
}
