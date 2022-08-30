package pathreader

import (
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"os"
)

//PathReader is a reader which will read from a shared path file
type PathReader struct {
	Path   string
	closed bool
	reader io.Reader
	closer io.Closer
}

func (pr *PathReader) SeekStart() error {
	return nil
}

func (pr *PathReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (pr *PathReader) Close() error {
	fmt.Println("path reader close")
	if !pr.closed {
		pr.closed = true

		if pr.closer != nil {
			return pr.closer.Close()
		}
	}
	return nil
}

func (pr *PathReader) Read(p []byte) (n int, err error) {
	if pr.closed {
		return 0, xerrors.Errorf("file reader closed")
	}

	if pr.reader == nil {
		fmt.Println("path reader read")
		fd, err := os.Open(pr.Path)

		if err != nil {
			pr.closed = true
			return 0, err
		}

		pr.reader = fd
		pr.closer = fd
	}

	return pr.reader.Read(p)
}

var _ io.ReadCloser = &PathReader{}
