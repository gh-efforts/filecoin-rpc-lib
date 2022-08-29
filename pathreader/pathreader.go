package pathreader

import (
	"fmt"
	carv2 "github.com/ipld/go-car/v2"
	"golang.org/x/xerrors"
	"io"
)

//PathReader is a reader which will read from a shared path file
type PathReader struct {
	Path string

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
	pr.Path = ""
	if pr.closer != nil {
		return pr.closer.Close()
	}
	return nil
}

func (pr *PathReader) Read(p []byte) (n int, err error) {
	if pr.reader == nil {
		fmt.Println("path reader read")
		rd, err := carv2.OpenReader(pr.Path)
		if err != nil {
			return 0, err
		}
		sr, err := rd.DataReader()
		if err != nil {
			return 0, err
		}

		// mark the reader as reading
		pr.Path = ""
		pr.reader = sr
		pr.closer = rd
	}
	if pr.reader == nil {
		return 0, xerrors.Errorf("file reader closed")
	}

	return pr.reader.Read(p)
}

var _ io.ReadCloser = &PathReader{}
