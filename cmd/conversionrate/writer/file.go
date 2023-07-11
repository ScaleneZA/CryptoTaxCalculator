package writer

import (
	"fmt"
	"io"
	"os"
)

const destination = "cmd/conversionrate/data"

type FileWriter struct {
	Filename string
}

func (w FileWriter) Write(reader io.Reader) error {
	fmt.Println(os.Getwd())
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
