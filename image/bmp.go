package image

import (
	"bytes"
	"encoding/binary"
)

type bmp struct{}

func init() {
	Register(bmp{})
}

func (bm bmp) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 6)

	// start: "BM"
	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}
	if !bytes.Equal(b[:2], []byte("BM")) {
		return
	}

	want := int64(binary.LittleEndian.Uint32(b[2:]))
	return ra.Size() >= want, nil
}

func (bm bmp) Exts() []string {
	return []string{".bmp"}
}
