package sharedtypes

type Pair struct {
	Currency1 string
	Currency2 string
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
