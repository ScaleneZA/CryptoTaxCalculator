package writer

import (
	"io"
)

type SQLLiteWriter struct {
	Filename string
}

func (w SQLLiteWriter) Write(reader io.Reader) error {
	// TODO: work
	return nil
}
