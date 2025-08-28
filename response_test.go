package bnm_test

import (
	"testing"

	"github.com/OsoianMarcel/bnm-go/v2"
)

func TestResponse_FindByCode(t *testing.T) {
	currency := bnm.Currency{Code: "USD", Name: "US Dollar", Value: 18.1434}
	response := bnm.Response{
		Currencies: []bnm.Currency{
			{Code: "EUR", Name: "Euro", Value: 20},
			currency,
		},
	}

	t.Run("existing currency", func(t *testing.T) {
		result, ok := response.FindByCode("USD")
		if !ok {
			t.Fatalf("currency USD not found")
		}
		if result != currency {
			t.Errorf("incorrect currency returned, expected: %+v, got: %+v", currency, result)
		}
	})

	t.Run("non-existent currency", func(t *testing.T) {
		if _, ok := response.FindByCode("INEXISTENT_CURRENCY"); ok {
			t.Errorf("inexistent currency found")
		}
	})
}
