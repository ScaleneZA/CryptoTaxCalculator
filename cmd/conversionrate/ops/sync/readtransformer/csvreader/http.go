package csvreader

import (
	"fmt"
	"net/http"
)

type HTTPCSVReader struct {
	Location string
	SkipRows int
}

func (r HTTPCSVReader) Read() ([][]string, error) {
	resp, err := http.Get(r.Location)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return readCSVFile(resp.Body, r.SkipRows)
}
