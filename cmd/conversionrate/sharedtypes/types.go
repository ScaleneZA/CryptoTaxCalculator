package sharedtypes

type Pair struct {
	FromCurrency string
	ToCurrency   string
}

type MarketSlice struct {
	Timestamp int
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
)

func AllPairs() []Pair {
	return []Pair{
		PairUSDBTC,
	}
}
