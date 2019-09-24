package main

import (
	"os"
)

type ReaderAt struct {
	file *os.File
}

func (r *ReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		st, err := r.file.Stat()
		if err != nil {
			return 0, err
		}
		off = st.Size() + off
	}
	return r.file.ReadAt(p, off)
}
