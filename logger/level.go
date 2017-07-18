package logger

import (
	"bytes"
	"io"
	"sync"
)

// LogLevel is
type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger
// that will filter out log message
type LevelFilter struct {
	Levels   []LogLevel
	MinLevel LogLevel
	Writer   io.Writer

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	f.once.Do(f.init)

	var level LogLevel
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = LogLevel(line[x+1 : x+y])
		}
	}

	_, ok := f.badLevels[level]
	return !ok
}

// Write is
func (f *LevelFilter) Write(p []byte) (n int, err error) {
	if !f.Check(p) {
		return len(p), nil
	}

	return f.Writer.Write(p)
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	badLevels := make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if level == f.MinLevel {
			break
		}
		badLevels[level] = struct{}{}
	}
	f.badLevels = badLevels
}
