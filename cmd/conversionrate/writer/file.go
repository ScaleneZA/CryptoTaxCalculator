package writer

import (
	"io"
	"os"
)

const destination = "../data"

type FileWriter struct {
	Filename string
}

func (w FileWriter) Write(reader io.Reader) error {
	out, err := os.Create(destination + "/" + w.Filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	return nil
}
