package conversionrate

func MarketValue(currencyFrom, currencyTo string, timestamp int) {

}

func SyncAll() error {
	for _, s := range Pairs {
		err := s.syncer.sync()
		if err != nil {
			return err
		}
	}

	return nil
}
