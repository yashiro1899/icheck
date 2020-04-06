package image

import "bytes"

type webp struct{}

func init() {
	Register(webp{})
}

func (p webp) Check(ra ReaderAt) (result bool, err error) {
	b := make([]byte, 20)

	// end: <?xpacket end='w'?>
	_, err = ra.ReadAt(b, -20)
	if err != nil {
		return
	}

	want := []byte("<?xpacket end='w'?>")
	return bytes.Equal(b[1:], want) || bytes.Equal(b[:19], want), nil
}

func (p webp) Exts() []string {
	return []string{".webp"}
}
