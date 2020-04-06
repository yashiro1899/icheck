package image

import "bytes"

type gif struct{}

func init() {
	Register(gif{})
}

func (g gif) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 2)

	// end: 003b
	_, err = ra.ReadAt(b, -2)
	if err != nil {
		return
	}
	if bytes.Equal(b, []byte("\x00\x3B")) {
		return true, nil
	}

	b = make([]byte, ra.Size()/2)
	_, err = ra.ReadAt(b, ra.Size()/-2)
	if err != nil {
		return
	}
	for i := len(b); i > 1; i-- {
		if bytes.Equal(b[i-2:i], []byte("\x00\x3B")) {
			return true, nil
		}
	}
	return
}

func (g gif) Exts() []string {
	return []string{".gif"}
}
