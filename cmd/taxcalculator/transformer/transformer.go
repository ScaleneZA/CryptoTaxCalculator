package transformer

import (
    "errors"
    "math"
    "sort"
    "strconv"

    "github.com/ScaleneZA/CryptoTaxCalculator/cmd/taxcalculator/sharedtypes"
    "github.com/xuri/excelize/v2"
)

func Transform(filename string, typ TransformType) ([]sharedtypes.Transaction, error) {
    rows, err := readFile(filename)
    if err != nil {
        return nil, err
    }

    var ts []sharedtypes.Transaction
    switch typ {
    case TransformTypeTest:
        ts, err = transformTest(rows)
        if err != nil {
            return nil, err
        }
    default:
        return nil, errors.New("invalid transaction type")
    }

    sort.Slice(ts, func(i, j int) bool {
        return ts[i].Timestamp < ts[j].Timestamp
    })

    return ts, nil
}

// TODO: Extract common code by using a mapping for the columns.
func transformTest(rows [][]string) ([]sharedtypes.Transaction, error) {
    var ts []sharedtypes.Transaction
    for i, r := range rows {
        amount, err := strconv.ParseFloat(r[2], 64)
        if err != nil {
            // Skip potential first row headers
            if i == 0 {
                continue
            }
            return nil, err
        }

        timestamp, err := strconv.Atoi(r[3])
        if err != nil {
            panic(err)
            return nil, err
        }

        typ := sharedtypes.TypeBuy
        if amount < 0 {
            amount = math.Abs(amount)
            typ = sharedtypes.TypeSell
        }

        wholePrice, err := strconv.ParseFloat(r[4], 64)
        if err != nil {
            return nil, err
        }

        ts = append(ts, sharedtypes.Transaction{
            Typ:               typ,
            Amount:            amount,
            Timestamp:         int64(timestamp),
            WholePriceAtPoint: wholePrice,
        })
    }

    return ts, nil
}

// TODO: Pull this out to its own package perhaps?
func readFile(filename string) ([][]string, error) {
    f, err := excelize.OpenFile(filename)
    if err != nil {
        return nil, err
    }
    defer func() {
        // Close the spreadsheet.
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()
    rows, err := f.GetRows("Sheet1")
    if err != nil {
        return nil, err
    }

    return rows, nil
}
