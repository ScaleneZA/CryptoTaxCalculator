package sharedtypes

type TransactionType int

const (
	TypeUnknown TransactionType = 0

	// TypeSendInternal and TypeReceiveInternal do not affect the tally.
	TypeSendInternal    TransactionType = 1
	TypeReceiveInternal TransactionType = 2

	TypeBuy             TransactionType = 3
	TypeSell            TransactionType = 4
	TypeSendExternal    TransactionType = 5
	TypeReceiveExternal TransactionType = 6
	TypeFee             TransactionType = 7

	typeSentinel TransactionType = 8
)

func ValidTransactionTypes() []TransactionType {
	var types []TransactionType
	for i := int(TypeUnknown) + 1; i < int(typeSentinel); i++ {
		types = append(types, TransactionType(i))
	}

	return types
}

var transactionTypeStrings = map[TransactionType]string{
	TypeUnknown:         "Unknown",
	TypeBuy:             "Buy",
	TypeSell:            "Sell",
	TypeSendExternal:    "Send (external)",
	TypeReceiveExternal: "Receive (external)",
	TypeFee:             "Fee",
	TypeSendInternal:    "Transfer (send)",
	TypeReceiveInternal: "Transfer (receieve)",
}

func (tt TransactionType) String() string {
	str, ok := transactionTypeStrings[tt]
	if ok {
		return str
	}

	return "Unknown"
}

func (tt TransactionType) Int() int {
	return int(tt)
}

func (tt TransactionType) ShouldIncreaseTally() bool {
	return tt == TypeBuy || tt == TypeReceiveExternal
}

// Reminder: Decreasing the tally does not always affect tax.
func (tt TransactionType) ShouldDecreaseTally() bool {
	return tt == TypeSell || tt == TypeSendExternal || tt == TypeFee
}

type Transaction struct {
	Currency      string
	DetectedType  TransactionType
	OverridedType TransactionType
	Amount        float64
	Timestamp     int64
	// TODO: Make this dynamic if possible
	WholePriceAtPoint float64
}

func (t Transaction) FinalType() TransactionType {
	if t.OverridedType != TypeUnknown {
		return t.OverridedType
	}

	return t.DetectedType
}
