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
	PairUSDBTC = Pair{
		FromCurrency: "USD",
		ToCurrency:   "BTC",
	}
	PairZARUSD = Pair{
		FromCurrency: "ZAR",
		ToCurrency:   "USD",
	}
)

func AllPairs() []Pair {
	return []Pair{
		PairUSDBTC,
		PairZARUSD,
	}
}
