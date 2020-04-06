package image

import (
	"bytes"
)

type png struct{}

func init() {
	Register(png{})
}

func (p png) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 5)

	// end: 44ae 4260 82
	_, err = ra.ReadAt(b, -5)
	if err != nil {
		return
	}
	return bytes.Equal(b, []byte("\x44\xAE\x42\x60\x82")), nil
}

func (p png) Exts() []string {
	return []string{".png"}
}
