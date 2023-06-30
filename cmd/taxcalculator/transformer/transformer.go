package transformer

import (
    "errors"
    "fmt"
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
    headerCount := 0
    for i, r := range rows {
        var t sharedtypes.Transaction
        switch typ {
        case TransformTypeTest:
            t, err = transformTestRow(r)
            if err != nil {
                if headerCount < 1 {
                    fmt.Println(fmt.Sprintf("Skipping row %d may be header", i))
                    headerCount++
                    continue
                }
                return nil, err
            }
        default:
            return nil, errors.New("invalid transaction type")
        }

        ts = append(ts, t)
    }

    sort.Slice(ts, func(i, j int) bool {
        return ts[i].Timestamp < ts[j].Timestamp
    })

    return ts, nil
}

// TODO: Extract common code by using a mapping for the columns.
func transformTestRow(r []string) (sharedtypes.Transaction, error) {
    amount, err := strconv.ParseFloat(r[2], 64)
    if err != nil {
        return sharedtypes.Transaction{}, err
    }

    timestamp, err := strconv.Atoi(r[3])
    if err != nil {
        return sharedtypes.Transaction{}, err
    }

    typ := sharedtypes.TypeBuy
    if amount < 0 {
        amount = math.Abs(amount)
        typ = sharedtypes.TypeSell
    }

    wholePrice, err := strconv.ParseFloat(r[4], 64)
    if err != nil {
        return sharedtypes.Transaction{}, err
    }

    return sharedtypes.Transaction{
        Typ:               typ,
        Amount:            amount,
        Timestamp:         int64(timestamp),
        WholePriceAtPoint: wholePrice,
    }, nil
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
