package conversionrate

type Pair struct {
	currency1 string
	currency2 string
	syncer    syncer
}
