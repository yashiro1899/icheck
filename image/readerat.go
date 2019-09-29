package image

import "os"

type sizeReaderAt struct {
	size int64
	fp   *os.File
}

func (ra sizeReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		off = ra.size + off
	}
	return ra.fp.ReadAt(p, off)
}

func (ra sizeReaderAt) Size() int64 {
	return ra.size
}

func (ra sizeReaderAt) Close() error {
	return ra.fp.Close()
}

func NewReaderAt(name string) (ReaderAt, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	return sizeReaderAt{size: fi.Size(), fp: fp}, nil
}
