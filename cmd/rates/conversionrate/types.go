package conversionrate

import "fmt"

type Pair struct {
	FromCurrency string
	ToCurrency   string
}

func (p Pair) String() string {
	return fmt.Sprintf("%s/%s", p.FromCurrency, p.ToCurrency)
}

type MarketSlice struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
}

type MarketPair struct {
	Pair
	MarketSlice
}

var (
	PairZARUSD = Pair{
		FromCurrency: "ZAR",
		ToCurrency:   "USD",
	}
	PairUSDBTC = Pair{
		FromCurrency: "USD",
		ToCurrency:   "BTC",
	}
	PairUSDETH = Pair{
		FromCurrency: "USD",
		ToCurrency:   "ETH",
	}
	PairUSDLTC = Pair{
		FromCurrency: "USD",
		ToCurrency:   "LTC",
	}
	PairUSDBCH = Pair{
		FromCurrency: "USD",
		ToCurrency:   "BCH",
	}
	PairUSDBAT = Pair{
		FromCurrency: "USD",
		ToCurrency:   "BAT",
	}
	PairUSDLINK = Pair{
		FromCurrency: "USD",
		ToCurrency:   "LINK",
	}
)

func AllPairs() []Pair {
	return []Pair{
		PairZARUSD,
		PairUSDBTC,
		PairUSDETH,
		PairUSDLTC,
		PairUSDBCH,
		PairUSDBAT,
		PairUSDLINK,
	}
}
