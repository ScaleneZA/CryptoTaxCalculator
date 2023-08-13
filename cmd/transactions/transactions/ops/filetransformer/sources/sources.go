package sources

import "github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"

var currencyMap = map[string]string{
	"XBT":  "BTC",
	"XXBT": "BTC",
	"XETH": "ETH",
	"XLTC": "LTC",
	"XXDG": "DOGE",
	"XICN": "ICN",
	"XETC": "ETC",
	"XXRP": "XRP",
	"ZUSD": "USD",
}

type Source interface {
	// TransformRow returns a slice of transactions because some sources contain the fees in the same row
	TransformRow(row []string) ([]transactions.Transaction, error)
}
