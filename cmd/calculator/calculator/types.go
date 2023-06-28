package calculator

type TransactionType int
const (
	TypeBuy = 0
	TypeSell	= 1
)

type Transaction struct {
	Typ TransactionType
	Amount float64
	Timestamp int64
	// TODO: Make this dynamic if possible
	WholePriceAtPoint float64
}
