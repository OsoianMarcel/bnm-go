package bnm_test

import (
	"testing"
	"time"

	"github.com/OsoianMarcel/bnm-go"
)

// Get specific date for testing
func getSpecificDate() (time.Time, error) {
	return time.Parse("2006-Jan-02", "2017-Aug-05")
}

// Get specific query for testing
func getSpecificQuery() (bnm.Query, error) {
	date, err := getSpecificDate()
	if err != nil {
		return bnm.Query{}, err
	}

	return bnm.NewQuery("ro", date), nil
}

// Test method DateToString()
func TestQuery_DateToString(t *testing.T) {
	query, err := getSpecificQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "05.08.2017"
	result := query.DateToString()

	if result != expected {
		t.Errorf("incorrect date, expected: %s, result: %s", expected, result)
	}
}

// Test method GenerateURI()
func TestQuery_GenerateURI(t *testing.T) {
	query, err := getSpecificQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "http://www.bnm.md/ro/official_exchange_rates?get_xml=1&date=05.08.2017"
	result := query.GenerateURI()
	if result != expected {
		t.Errorf("incorrect URI, expected: %s, result: %s", expected, result)
	}
}

// Test method GetID()
func TestQuery_GetID(t *testing.T) {
	query, err := getSpecificQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "ro_05.08.2017"
	result := query.GetID()
	if result != expected {
		t.Errorf("incorrect id, expected: %s, result: %s", expected, result)
	}
}
