package image

import (
	"io"
	"mime"
	"net/http"
	"strings"
)

type Checker interface {
	Check(ReaderAt) (bool, error)
	Exts() []string
}

type ReaderAt interface {
	io.Closer
	io.ReaderAt
	Size() int64
}

var m = make(map[string]Checker)

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

func Sniff(ra ReaderAt) (ch Checker, err error) {
	b := make([]byte, 32)

	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}

	ct := http.DetectContentType(b)
	exts, err := mime.ExtensionsByType(ct)
	if err != nil {
		return
	}

	for _, ext := range exts {
		ch = Get(ext)
		if ch != nil {
			return
		}
	}
	return
}
