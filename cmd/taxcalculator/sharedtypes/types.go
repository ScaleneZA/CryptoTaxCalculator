package sharedtypes

type TransactionType int

const (
	TypeUnknown         TransactionType = 0
	TypeBuy             TransactionType = 1
	TypeSell            TransactionType = 2
	TypeSendExternal    TransactionType = 3
	TypeReceiveExternal TransactionType = 4
	TypeFee             TransactionType = 5

	// TypeSendInternal and TypeReceiveInternal do not affect the tally.
	TypeSendInternal    TransactionType = 6
	TypeReceiveInternal TransactionType = 7
)

func (tt TransactionType) ShouldIncreaseTally() bool {
	return tt == TypeBuy || tt == TypeReceiveExternal
}

// Reminder: Decreasing the tally does not always affect tax.
func (tt TransactionType) ShouldDecreaseTally() bool {
	return tt == TypeSell || tt == TypeSendExternal || tt == TypeFee
}

type Transaction struct {
	Currency  string
	Typ       TransactionType
	Amount    float64
	Timestamp int64
	// TODO: Make this dynamic if possible
	WholePriceAtPoint float64
}
