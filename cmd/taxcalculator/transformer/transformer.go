package transformer

import (
    "encoding/csv"
    "errors"
    "fmt"
    "os"
    "sort"

    "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
    "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/transformer/sources"
)

func Transform(filename string, typ TransformType) ([]sharedtypes.Transaction, error) {
    rows, err := readFile(filename)
    if err != nil {
        return nil, err
    }

    src, err := sourceFromType(typ)
    if err != nil {
        return nil, err
    }

    var ts []sharedtypes.Transaction
    headerCount := 0
    for i, r := range rows {
        t, err := src.TransformRow(r)
        if err != nil {
            if headerCount < 1 {
                fmt.Println(fmt.Sprintf("Skipping row %d may be header", i))
                headerCount++
                continue
            }
            return nil, err
        }

        ts = append(ts, t)
    }

    sort.Slice(ts, func(i, j int) bool {
        return ts[i].Timestamp < ts[j].Timestamp
    })

    return ts, nil
}

func sourceFromType(typ TransformType) (sources.Source, error) {
    var src sources.Source
    switch typ {
    case TransformTypeBasic:
        src = sources.BasicSource{}
    case TransformTypeLuno:
        src = sources.LunoSource{}
    default:
        return nil, errors.New("invalid source")
    }

    return src, nil
}

// TODO: Pull this out to its own package perhaps?
func readFile(filename string) ([][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    rows, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    return rows, nil
}
