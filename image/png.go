package image

import (
	"io"
)

type Png struct{}

func init() {
	Register(new(Png))
}

func (p *Png) Check(ra io.ReaderAt) (result bool, err error) {
	b := make([]byte, 8)

	// start: 8950 4e47 0d0a 1a0a
	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}
	if !(b[0] == 0x89 &&
		b[1] == 0x50 &&
		b[2] == 0x4e &&
		b[3] == 0x47 &&
		b[4] == 0x0d &&
		b[5] == 0x0a &&
		b[6] == 0x1a &&
		b[7] == 0x0a) {
		return
	}

	// end: 44ae 4260 82
	_, err = ra.ReadAt(b, -8)
	if err != nil {
		return
	}
	if b[3] == 0x44 &&
		b[4] == 0xae &&
		b[5] == 0x42 &&
		b[6] == 0x60 &&
		b[7] == 0x82 {
		return true, nil
	}
	return
}

func (p *Png) Exts() []string {
	return []string{".png"}
}
