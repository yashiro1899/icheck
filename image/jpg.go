package image

import (
	"io"
)

type Jpg struct{}

func init() {
	Register(new(Jpg))
}

func (j *Jpg) Check(ra io.ReaderAt) error {
	b := make([]byte, 2)

	// start: ffd8
	_, err := ra.ReadAt(b, 0)
	if err != nil {
		return err
	}
	if !(b[0] == 0xff && b[1] == 0xd8) {
		return Incomplete
	}

	// end: ffd9
	for i := -2; i > -64; i-- {
		_, err := ra.ReadAt(b, int64(i))
		if err != nil {
			return err
		}
		if b[0] == 0xff && b[1] == 0xd9 {
			return nil
		}
	}
	return Incomplete
}

func (j *Jpg) Exts() []string {
	return []string{".jpg", ".jpeg"}
}
