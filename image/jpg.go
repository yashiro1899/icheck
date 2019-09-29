package image

import (
	"bytes"
)

type jpg struct{}

func init() {
	Register(jpg{})
}

func (j jpg) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 2)

	// start: ffd8
	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}
	if !bytes.Equal(b, []byte("\xFF\xD8")) {
		return
	}

	// end: ffd9
	for i := -2; i > -64; i-- {
		_, err = ra.ReadAt(b, int64(i))
		if err != nil {
			return
		}
		if bytes.Equal(b, []byte("\xFF\xD9")) {
			return true, nil
		}
	}
	return
}

func (j jpg) Exts() []string {
	return []string{".jpg", ".jpeg"}
}
