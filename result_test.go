package bnm_test

import (
	"testing"
	"github.com/OsoianMarcel/bnm-go"
)

// Get Result struct filled for testing
func getSpecificResult() bnm.Result {
	rates := make([]bnm.Rate, 0, 3)
	rates = append(
		rates,
		bnm.Rate{Code: "EUR", Name: "Euro", Value: 20.5964},
		bnm.Rate{Code: "USD", Name: "US Dollar", Value: 18.1434},
		bnm.Rate{Code: "RON", Name: "Romanian Leu", Value: 4.4931},
	)

	return bnm.Result{Rates: rates}
}

// Test Result method FindByCode()
func TestResult_FindByCode(t *testing.T) {
	result := getSpecificResult()

	currencyName := "USD"
	expected := bnm.Rate{Code: "USD", Name: "US Dollar", Value: 18.1434}
	if result, ok := result.FindByCode(currencyName); ok {
		if result != expected {
			t.Errorf("incorrect rate value, expected: %+v, result: %+v", expected, result)
		}
	} else {
		t.Errorf("currency %s not found", currencyName)
	}

	if _, ok := result.FindByCode("INEXISTENT_CURRENCY"); ok {
		t.Error("inexistent currency found")
	}
}