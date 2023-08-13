package transactions

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
	TypeInterest        TransactionType = 8
	TypeAirdrop         TransactionType = 9

	typeSentinel TransactionType = 10
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
	return tt == TypeBuy || tt == TypeReceiveExternal || tt == TypeInterest || tt == TypeAirdrop
}

// Reminder: Decreasing the tally does not always affect tax.
func (tt TransactionType) ShouldDecreaseTally() bool {
	return tt == TypeSell || tt == TypeSendExternal || tt == TypeFee
}

// ShouldCheck supplies types that need double-checking by the user.
func (tt TransactionType) ShouldCheck() bool {
	return tt == TypeSendInternal || tt == TypeReceiveInternal
}

type FiatPrice struct {
	Fiat  string
	Price float64
}

type Transaction struct {
	UID           string
	Transformer   TransformType
	Currency      string
	DetectedType  TransactionType
	OverridedType TransactionType
	Amount        float64
	Timestamp     int64

	WholePriceAtPoint FiatPrice
}

func (t Transaction) FinalType() TransactionType {
	if t.OverridedType != TypeUnknown {
		return t.OverridedType
	}

	return t.DetectedType
}

type TransformType int

const (
	TransformTypeUnknown  = 0
	TransformTypeBasic    = 1
	TransformTypeLuno     = 2
	TransformTypeBinance  = 3
	TransformTypeCoinomi  = 4
	TransformTypeKraken   = 5
	transformTypeSentinel = 6
)

func SelectableTransformTypes() []TransformType {
	return []TransformType{
		TransformTypeBasic,
		TransformTypeBinance,
		TransformTypeCoinomi,
		TransformTypeLuno,
		TransformTypeKraken,
	}
}

var transformTypeStrings = map[TransformType]string{
	TransformTypeUnknown: "Unknown",
	TransformTypeBasic:   "Basic (Don't use)",
	TransformTypeLuno:    "Luno",
	TransformTypeBinance: "Binance",
	TransformTypeCoinomi: "Coinomi",
	TransformTypeKraken:  "Kraken",
}

func (tt TransformType) String() string {
	str, ok := transformTypeStrings[tt]
	if ok {
		return str
	}

	return "Unknown"
}

func (tt TransformType) Int() int {
	return int(tt)
}
