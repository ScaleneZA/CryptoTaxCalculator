package sharedtypes

type TransactionType int

const (
	TypeBuy             TransactionType = 0
	TypeSell            TransactionType = 1
	TypeSendExternal    TransactionType = 2
	TypeReceiveExternal TransactionType = 3

	// TypeSendInternal and TypeSendExternal do not affect the tally.
	TypeSendInternal    TransactionType = 4
	TypeReceiveInternal TransactionType = 5
)

func (tt TransactionType) ShouldIncreaseTally() bool {
	return tt == TypeBuy || tt == TypeReceiveExternal
}

// Reminder: Decreasing the tally does not always affect tax.
func (tt TransactionType) ShouldDecreaseTally() bool {
	return tt == TypeSell || tt == TypeSendExternal
}

type Transaction struct {
	Currency  string
	Typ       TransactionType
	Amount    float64
	Timestamp int64
	// TODO: Make this dynamic if possible
	WholePriceAtPoint float64
}
