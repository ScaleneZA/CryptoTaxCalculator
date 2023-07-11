package writer

import "io"

type Writer interface {
	Write(reader io.Reader) error
}
