package calculator

type YearEndTotal struct {
	Year     int
	Gains    []Gain
	Balances []Balance
}

type Gain struct {
	Asset    string
	Amount   float64
	Costs    float64
	Proceeds float64
}

type Balance struct {
	Asset  string
	Amount float64
	// TODO: Add Costs and MarketValue
}
