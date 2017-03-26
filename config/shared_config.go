package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type dataSource interface {
	ReadCloser() (io.ReadCloser, error)
}

type sourceFile struct {
	name string
}

func (s sourceFile) ReadCloser() (io.ReadCloser, error) {
	return os.Open(s.name)
}

// File represents a combination of a json file in memory.
type File struct {
	dataSource dataSource
	v          interface{}
}

func newFile(source dataSource, v interface{}) *File {
	return &File{
		dataSource: source,
		v:          v,
	}
}

// Load returns a new File pointer
func Load(name string, v interface{}) (*File, error) {
	source, err := parseDataSource(name)
	if err != nil {
		return nil, err
	}

	f := newFile(source, v)
	if err := f.Reload(); err != nil {
		return nil, err
	}

	return f, nil
}

// Reload reloads and parses data source
func (f *File) Reload() (err error) {
	if err := f.reload(f.dataSource); err != nil {
		return err
	}

	return nil
}

func (f *File) reload(s dataSource) error {
	r, err := s.ReadCloser()
	if err != nil {
		return err
	}
	defer r.Close()

	return f.parse(r)
}

func parseDataSource(filename string) (dataSource, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}

	if homeDir == "" {
		return nil, errors.New("user home directory not found")
	}

	path := filepath.Join(homeDir, ".feiniubus", filename)
	return sourceFile{path}, nil
}
