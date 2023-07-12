package writer

import (
	"encoding/csv"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/conversionrate/sync/sharedtypes"
	"os"
	"strconv"
)

const destination = "cmd/conversionrate/sync/data"

type FileWriter struct {
	Filename string
}

func (w FileWriter) Write(mps []sharedtypes.MarketSlice) error {
	out, err := os.Create(destination + "/" + w.Filename)
	if err != nil {
		return err
	}
	defer out.Close()

	var records [][]string
	for _, mp := range mps {
		records = append(records, []string{
			strconv.Itoa(mp.Timestamp),
			strconv.FormatFloat(mp.Open, 'f', -1, 64),
			strconv.FormatFloat(mp.High, 'f', -1, 64),
			strconv.FormatFloat(mp.Low, 'f', -1, 64),
			strconv.FormatFloat(mp.Close, 'f', -1, 64),
		})
	}

	csvW := csv.NewWriter(out)

	err = csvW.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}
