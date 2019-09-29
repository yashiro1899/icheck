package image

import (
	"bytes"
)

type Png struct{}

func init() {
	Register(new(Png))
}

func (p *Png) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 8)

	// start: 8950 4e47 0d0a 1a0a
	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}
	if !bytes.Equal(b, []byte("\x89PNG\x0D\x0A\x1A\x0A")) {
		return
	}

	// end: 44ae 4260 82
	b = b[:5]
	_, err = ra.ReadAt(b, -5)
	if err != nil {
		return
	}
	return bytes.Equal(b, []byte("\x44\xAE\x42\x60\x82")), nil
}

func (p *Png) Exts() []string {
	return []string{".png"}
}
