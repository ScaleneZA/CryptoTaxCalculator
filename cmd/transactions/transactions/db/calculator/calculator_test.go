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

func TestListByTypeByUid(t *testing.T) {
	dbc := db.ConnectForTesting()

	_, err := calculator.Upsert(dbc, "1234", transactions.TypeAirdrop)
	jtest.RequireNil(t, err)
	_, err = calculator.Upsert(dbc, "1235", transactions.TypeFee)
	jtest.RequireNil(t, err)
	_, err = calculator.Upsert(dbc, "1236", transactions.TypeSell)
	jtest.RequireNil(t, err)

	tests := []struct {
		name     string
		uids     []string
		expected []transactions.OverrideType
	}{
		{
			name: "Fetch only supplied UIDs",
			uids: []string{"1234", "1236"},
			expected: []transactions.OverrideType{
				{
					UID:            "1234",
					OverriddenType: transactions.TypeAirdrop,
				},
				{
					UID:            "1236",
					OverriddenType: transactions.TypeSell,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := calculator.ListByTypeByUid(dbc, tt.uids)
			jtest.RequireNil(t, err)

			require.Equal(t, tt.expected, actual)
		})
	}
}
