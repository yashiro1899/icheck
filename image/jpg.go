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
	_, err = ra.ReadAt(b, -2)
	if err != nil {
		return
	}
	if bytes.Equal(b, []byte("\xFF\xD9")) {
		return true, nil
	}

	b = make([]byte, ra.Size()/2)
	_, err = ra.ReadAt(b, ra.Size()/-2)
	if err != nil {
		return
	}
	for i := len(b); i > 1; i-- {
		if bytes.Equal(b[i-2:i], []byte("\xFF\xD9")) {
			return true, nil
		}
	}
	return
}

func (j jpg) Exts() []string {
	return []string{".jpg", ".jpeg"}
}
