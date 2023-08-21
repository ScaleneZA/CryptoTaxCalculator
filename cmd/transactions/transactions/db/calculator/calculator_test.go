package calculator_test

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/db"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/db/calculator"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpsertAndLookup(t *testing.T) {
	dbc := db.ConnectForTesting()

	_, err := calculator.Upsert(dbc, "4321", transactions.TypeAirdrop)
	jtest.RequireNil(t, err)

	tests := []struct {
		name string
		uid  string
		typ  transactions.TransactionType
	}{
		{
			name: "New entry",
			uid:  "1234",
			typ:  transactions.TypeFee,
		},
		{
			name: "Existing Entry, updates type",
			uid:  "4321",
			typ:  transactions.TypeInterest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculator.Upsert(dbc, tt.uid, tt.typ)
			jtest.RequireNil(t, err)

			typ, err := calculator.LookupTypeByUid(dbc, tt.uid)
			jtest.RequireNil(t, err)
			require.Equal(t, tt.typ, typ)
		})
	}
}
