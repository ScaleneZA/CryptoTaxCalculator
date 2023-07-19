package csvreader

import (
	"os"
)

type LocalCSVReader struct {
	Location string
	SkipRows int
}

func (r LocalCSVReader) Read() ([][]string, error) {
	file, err := os.Open(r.Location)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readCSVFile(file, r.SkipRows)
}
