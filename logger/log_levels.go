package logger

import (
	"io/ioutil"
)

// NewLevelFilter returns a LevelFilter that is configured with the log
// levels that we use.
func NewLevelFilter() *LevelFilter {
	return &LevelFilter{
		Levels:   []LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERR"},
		MinLevel: "INFO",
		Writer:   ioutil.Discard,
	}
}

// ValidateLevelFilter verifies that the log levels within the filter
// are valid.
func ValidateLevelFilter(minLevel LogLevel, filter *LevelFilter) bool {
	for _, level := range filter.Levels {
		if level == minLevel {
			return true
		}
	}
	return false
}
