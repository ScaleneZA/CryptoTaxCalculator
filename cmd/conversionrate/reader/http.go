package reader

import (
	"fmt"
	"io"
	"net/http"
)

type HttpReader struct {
	Location string
}

func (r HttpReader) Read() (io.Reader, error) {
	resp, err := http.Get(r.Location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, nil
}
