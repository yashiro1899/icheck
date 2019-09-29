package image

import (
	"bytes"
)

type Jpg struct{}

func init() {
	Register(new(Jpg))
}

func (j *Jpg) Check(ra ReaderAt) (result bool, err error) {
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

func (j *Jpg) Exts() []string {
	return []string{".jpg", ".jpeg"}
}
