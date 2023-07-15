package csvreader

import (
	"bufio"
	"encoding/csv"
	"io"
)

type Reader interface {
	Read() ([][]string, error)
}

func readCSVFile(file io.Reader, skipRows int) ([][]string, error) {
	buf, ok := (file).(*bufio.Reader)
	if !ok {
		buf = bufio.NewReader(file)
	}

	for i := 0; i < skipRows; i++ {
		_, err := buf.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
	}

	reader := csv.NewReader(buf)
	return reader.ReadAll()
}
