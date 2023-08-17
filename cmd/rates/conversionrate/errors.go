package conversionrate

import (
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

var (
	ErrStoredRateExceedsThreshold = errors.New("closest timestamps of stored rates exceed threshold of 1 week", j.C("ERR_2241f733af924a9a"))
	ErrNoMarket                   = errors.New("no market found for pair within bounds", j.C("ERR_2241f733af924a9b"))
	ErrNoRatesFound               = errors.New("no rates found", j.C("ERR_2241f733af924a9c"))
	ErrPairSyncFailed             = errors.New("one or more pairs failed to sync", j.C("ERR_2241f733af924a9d"))
)
