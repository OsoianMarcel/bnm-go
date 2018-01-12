package bnm_test

import (
	"testing"
	"github.com/osoianmarcel/bnm"
	"time"
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

// Test method GenerateUri()
func TestQuery_GenerateUri(t *testing.T) {
	query, err := getSpecificQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "http://www.bnm.md/ro/official_exchange_rates?get_xml=1&date=05.08.2017"
	result := query.GenerateUri() 
	if result != expected {
		t.Errorf("incorrect URI, expected: %s, result: %s", expected, result)
	}
}

// Test method GetId()
func TestQuery_GetId(t *testing.T) {
	query, err := getSpecificQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "ro_05.08.2017"
	result := query.GetId() 
	if result != expected {
		t.Errorf("incorrect id, expected: %s, result: %s", expected, result)
	}
}