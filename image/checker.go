package image

import (
	"errors"
	"io"
	"strings"
)

type Checker interface {
	Check(io.ReaderAt) error
	Exts() []string
}

var (
	Incomplete = errors.New("incomplete")
	m          = make(map[string]Checker)
)

func Register(c Checker) {
	for _, ext := range c.Exts() {
		m[ext] = c
	}
}

func Get(ext string) Checker {
	if c, ok := m[strings.ToLower(ext)]; ok {
		return c
	}
	return nil
}
