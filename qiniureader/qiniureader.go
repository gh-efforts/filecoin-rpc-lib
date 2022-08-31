package qiniureader

import (
	"encoding/json"
	"github.com/qiniupd/qiniu-go-sdk/syncdata/operation"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"os"
)

type QiniuReader struct {
	Key string
	// filePath is optional
	filePath string
	closed   bool
	body     io.ReadCloser
}

func (reader *QiniuReader) SeekStart() error {
	return nil
}

func (reader *QiniuReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (reader *QiniuReader) Close() error {
	if !reader.closed {
		reader.closed = true

		if reader.body != nil {
			return reader.body.Close()
		}
	}
	return nil
}

func (reader *QiniuReader) Read(p []byte) (n int, err error) {
	if reader.closed {
		return 0, xerrors.Errorf("file reader closed")
	}

	if reader.body == nil {
		configPath := os.Getenv("QINIU_READER_CONFIG_PATH")
		configBytes, err := os.ReadFile(configPath)

		if err != nil {
			return 0, err
		}

		var config operation.Config

		err = json.Unmarshal(configBytes, &config)

		if err != nil {
			return 0, err
		}

		var resp *http.Response
		resp, err = operation.NewDownloader(&config).DownloadRaw(reader.Key, http.Header{})

		if err != nil {
			return 0, err
		}

		reader.body = resp.Body
	}

	return reader.body.Read(p)
}
