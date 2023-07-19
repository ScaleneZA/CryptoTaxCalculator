package conversionrate

type Client interface {
	ValueAtTime(from, to string, timestamp int64) (float64, error)
}
