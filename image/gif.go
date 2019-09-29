package image

import "bytes"

type gif struct{}

func init() {
	Register(gif{})
}

func (g gif) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 6)

	// start: "GIF87a" or "GIF89a"
	_, err = ra.ReadAt(b, 0)
	if err != nil {
		return
	}
	if !bytes.Equal(b, []byte("GIF87a")) && !bytes.Equal(b, []byte("GIF89a")) {
		return
	}

	// end: 003b
	b = b[:2]
	for i := int64(-2); i > ra.Size()/-2; i-- {
		_, err = ra.ReadAt(b, i)
		if err != nil {
			return
		}
		if bytes.Equal(b, []byte("\x00\x3B")) {
			return true, nil
		}
	}
	return
}

func (g gif) Exts() []string {
	return []string{".gif"}
}
