package bnm_test

import (
	"testing"
	"time"

	"github.com/OsoianMarcel/bnm-go/v2"
)

func getSpecificDate() time.Time {
	t, err := time.Parse("2006-Jan-02", "2017-Aug-05")
	if err != nil {
		panic(err)
	}

	return t
}

func getSpecificQuery() bnm.Query {
	return bnm.NewQuery(getSpecificDate(), bnm.LANG_EN)
}

// Test method RequestURL()
func TestQuery_RequestURL(t *testing.T) {
	query := getSpecificQuery()

	expected := "http://www.bnm.md/en/official_exchange_rates?get_xml=1&date=05.08.2017"
	result := query.RequestURL()
	if result != expected {
		t.Errorf("incorrect URL, expected: %s, result: %s", expected, result)
	}
}

// Test method GetID()
func TestQuery_ID(t *testing.T) {
	query := getSpecificQuery()

	expected := "en_05.08.2017"
	result := query.ID()
	if result != expected {
		t.Errorf("incorrect id, expected: %s, result: %s", expected, result)
	}
}
