package sharedtypes

type Pair struct {
	FromCurrency string
	ToCurrency   string
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
)

func AllPairs() []Pair {
	return []Pair{
		PairUSDBTC,
		PairZARUSD,
	}
}
