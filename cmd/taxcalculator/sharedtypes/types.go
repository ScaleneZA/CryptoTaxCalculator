package sharedtypes

type TransactionType int

const (
	TypeBuy  TransactionType = 0
	TypeSell TransactionType = 1
)

type Transaction struct {
	Currency  string
	Typ       TransactionType
	Amount    float64
	Timestamp int64
	// TODO: Make this dynamic if possible
	WholePriceAtPoint float64
}
