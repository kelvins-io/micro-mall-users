package convert

import (
	"bytes"
	"io"
)

func IoReaderToByte(reader io.Reader) *bytes.Buffer {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	return buf
}
