package image

import (
	"encoding/binary"
)

type webp struct{}

func init() {
	Register(webp{})
}

func (p webp) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 8)

	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}

	want := int64(binary.LittleEndian.Uint32(b[4:])) + 8
	return ra.Size() >= want, nil
}

func (p webp) Exts() []string {
	return []string{".webp"}
}
