package sources

import (
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

var (
	ErrSkipTransaction = errors.New("Transaction is not valid, consider skipping", j.C("ERR_1241f563af924a9c"))
)
