package taxcalculator

import (
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

var (
	ErrUnsupportedTranformType = errors.New("unsupported transaction source", j.C("ERR_1241f733af924a9a"))
	ErrInvalidTransactionOrder = errors.New("transactions have not been ordered correctly", j.C("ERR_1241f733af924a9b"))
)
