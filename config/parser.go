package config

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
)

type parser struct {
	buf   *bufio.Reader
	isEOF bool
}

func newParser(r io.Reader) *parser {
	return &parser{
		buf: bufio.NewReader(r),
	}
}

func (p *parser) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(p.buf)
}

func (f *File) parse(reader io.Reader) (err error) {
	p := newParser(reader)
	buf, err := p.ReadAll()
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &f.v)
	return err
}
