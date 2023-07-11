package sync

type Pair struct {
	currency1 string
	currency2 string
}

type MarketPair struct {
	pair      Pair
	timestamp int
	open      float64
	high      float64
	low       float64
	close     float64
}
