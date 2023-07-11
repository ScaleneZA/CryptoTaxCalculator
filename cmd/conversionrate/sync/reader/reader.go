package reader

import "io"

type Reader interface {
	Read() (io.ReadCloser, error)
}
